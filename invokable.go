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

type resultEnvelope struct {
	result []reflect.Value
	actIdx int
	maxIdx int
}

type resultFuture chan resultEnvelope

type Invokable interface {
	receive() []reflect.Value
}

type invokable struct {
	err error
	rf  resultFuture
	pr  *progressor
}

///////////////////////////////////////////////////////////////////////////////////////
// setError
///////////////////////////////////////////////////////////////////////////////////////
func (e *invokable) setError(fnv interface{}, err error) {

	if fnv == nil {
		e.err = err
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
	fmt.Fprintf(buffer, "%s: line %d -> %s '%s'\n\n", fileName, line, funcName, err.Error())

	e.err = errors.New(buffer.String())
}

///////////////////////////////////////////////////////////////////////////////////////
// sendError
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) sendError(fnv interface{}, idx int, err error) {
	// send dummy to end receiver and avoid goroutine deadlock
	i.send([]reflect.Value{}, idx, -1)
	i.setError(fnv, err)
}

///////////////////////////////////////////////////////////////////////////////////////
// send
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) send(out []reflect.Value, actIdx int, maxIdx int) {
	i.rf <- resultEnvelope{result: out, actIdx: actIdx, maxIdx: maxIdx}
}

///////////////////////////////////////////////////////////////////////////////////////
// receive
///////////////////////////////////////////////////////////////////////////////////////
func (i *invokable) receive() []reflect.Value {

	nInputs := 0
	rvd := [][]reflect.Value{}

	insert := func(e resultEnvelope) bool {
		for e.actIdx >= len(rvd) {
			rvd = append(rvd, []reflect.Value{})
		}

		rvd[e.actIdx] = e.result
		nInputs++

		//in case of composing error in.maxIdx == -1 to end loop
		return nInputs < e.maxIdx && e.maxIdx != -1
	}

	for insert(<-i.rf) {
	}

	//flatten received data
	data := []reflect.Value{}
	for _, arr := range rvd {
		data = append(data, arr...)
	}

	return data
}

///////////////////////////////////////////////////////////////////////////////////////
// invokeTargets
///////////////////////////////////////////////////////////////////////////////////////
func (p *invokable) invokeTargets(targets []reflect.Value, inputs []reflect.Value) {

	maxIdx := len(targets)
	for idx, target := range targets {
		switch t := target.Type(); t.Kind() {
		case reflect.Func:
			r := Resolver(p, target)
			nOut := r.InArgCount()

			var err error
			if !r.CanInvokeWithParams(inputs) {
				inputs, err = r.Resolve(inputs, func(resInput []reflect.Value) {
					p.send(target.Call(resInput), idx, maxIdx)
				})

				if err != nil {
					p.sendError(target, idx, err)
					publishError(err)
					return
				}

			} else {
				// extract the inputs we need and invoke the func
				actIn := inputs[:nOut]
				inputs = inputs[nOut:]
				p.send(target.Call(actIn), idx, maxIdx)
			}

		default:
			//send values directly
			p.send([]reflect.Value{target}, idx, maxIdx)
		}
	}

	if len(inputs) > 0 {
		err := fmt.Errorf("Unused inputs on target. %v\n", fmt.Sprint(inputs))
		publishError(err)
	}
}
