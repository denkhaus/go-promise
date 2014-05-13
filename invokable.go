package Q

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

type errorAware struct {
	err error
}

type resultEnvelope struct {
	result []reflect.Value
	actIdx int
	maxIdx int
}

type ParamFuture chan []reflect.Value

type ResultFuture chan resultEnvelope

type invokable struct {
	errorAware
	pf ParamFuture
	rf ResultFuture
}

///////////////////////////////////////////////////////////////////////////////////////
// setError
///////////////////////////////////////////////////////////////////////////////////////
func (e *errorAware) setError(fnv interface{}, err string) {

	if fnv == nil {
		e.err = errors.New(err)
		return
	}

	pc := fnv.(reflect.Value).Pointer()
	f := runtime.FuncForPC(pc)
	file, line := f.FileLine(pc)
	_, fileName := path.Split(file)
	funcNameParts := strings.Split(f.Name(), ".")
	funcNamePartsIdx := len(funcNameParts) - 1
	funcName := funcNameParts[funcNamePartsIdx]
	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "Q internal error ------------------------------------------------------------------------\n")
	fmt.Fprintf(buffer, "%s: line %d -> %s '%s'\n\n", fileName, line, funcName, err)

	e.err = errors.New(buffer.String())
}

///////////////////////////////////////////////////////////////////////////////////////
// send
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) send(out []reflect.Value) {
	i.pf <- out
}

///////////////////////////////////////////////////////////////////////////////////////
// receive
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) receive() (out []reflect.Value) {

	outR := false

	for !outR {
		select {
		case out = <-i.pf:
			outR = true
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////
// receive
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) sendError(fnv interface{}, err string) {
	// send dummy to avoid goroutine deadlock
	i.sendWithIndex([]reflect.Value{}, -1, -1)
	i.setError(fnv, err)
}

///////////////////////////////////////////////////////////////////////////////////////
// invoke
///////////////////////////////////////////////////////////////////////////////////////
func (p *invokable) invoke(fn interface{}, in []reflect.Value) {
	v := reflect.ValueOf(fn)
	t := v.Type()

	//check arguments count equal
	if len(in) != t.NumIn() {
		p.send([]reflect.Value{}) // send dummy to avoid goroutine deadlock
		p.setError(v, "Function argument count mismatch.")
		return
	}
	//check arguments types equal
	for idx, inVal := range in {
		if inVal.Type() != t.In(idx) {
			p.send([]reflect.Value{}) // send dummy to avoid goroutine deadlock
			p.setError(v, "Function argument type mismatch.")
			return
		}
	}

	p.send(v.Call(in))
}

///////////////////////////////////////////////////////////////////////////////////////
// sendWithIndex
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) sendWithIndex(out []reflect.Value, actIdx int, maxIdx int) {
	i.rf <- resultEnvelope{result: out, actIdx: actIdx, maxIdx: maxIdx}
}

///////////////////////////////////////////////////////////////////////////////////////
// receiveWithIndex
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) receiveWithIndex() []reflect.Value {

	nInputs := 0
	rvd := [][]reflect.Value{}

	insert := func(e resultEnvelope) {
		for e.actIdx >= len(rvd) {
			rvd = append(rvd, []reflect.Value{})
		}
		rvd[e.actIdx] = e.result
	}

	for {
		select {
		case in := <-i.rf:
			nInputs++
			insert(in)

			if nInputs == in.maxIdx || // all inputs received
				in.maxIdx == -1 { //in case of internal error in.maxIdx == -1 to end loop
				break
			}
		}
	}

	//flatten received data
	data := []reflect.Value{}
	for _, arr := range rvd {
		data = append(data, arr...)
	}

	return data
}

///////////////////////////////////////////////////////////////////////////////////////
// invoke
///////////////////////////////////////////////////////////////////////////////////////
func (p *invokable) invokeFunc(fn reflect.Value, in []reflect.Value, targetIdx int, maxIdx int) {
	t := fn.Type()

	//check arguments types equal
	for idx, inVal := range in {
		if inVal.Type() != t.In(idx) {
			p.sendError(fn, "Function argument type mismatch.")
			return
		}
	}

	p.sendWithIndex(fn.Call(in), targetIdx, maxIdx)
}

///////////////////////////////////////////////////////////////////////////////////////
// invoke
///////////////////////////////////////////////////////////////////////////////////////
func (p *invokable) resolveQ(in []reflect.Value) []reflect.Value {

	out := []reflect.Value{}
	for _, inVal := range in {
		t := inVal.Type()

		//is input promis or deferred
		if t == PromiseType || t == DeferredType {
			v := inVal.MethodByName("receiveWithIndex")
			res := v.Call([]reflect.Value{})
			out = append(out, res...)
		} else {
			out = append(out, inVal)
		}
	}

	return out
}

///////////////////////////////////////////////////////////////////////////////////////
// invokeTargets
///////////////////////////////////////////////////////////////////////////////////////
func (p *invokable) invokeTargets(targets []reflect.Value, in []reflect.Value) {

	inputs := p.resolveQ(in)
	maxIdx := len(targets)

	for idx, target := range targets {
		t := target.Type()

		switch t.Kind() {
		case reflect.Func:
			nFnInpts := t.NumIn()

			//check we have enough func inputs
			if len(inputs) < nFnInpts {
				p.sendError(target, "Function argument count mismatch.")
				return
			}

			// extract the inputs we need and invoke the func
			actIn := inputs[:nFnInpts]
			inputs = inputs[nFnInpts:]
			p.invokeFunc(target, actIn, idx, maxIdx)

		default:
			//send values directly
			p.sendWithIndex([]reflect.Value{target}, idx, maxIdx)
		}
	}

	if len(inputs) > 0 {
		panic("invokeAll:: we got invoke inputs leftovers")
	}
}
