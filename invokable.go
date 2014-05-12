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

type ParamFuture chan []reflect.Value

type invokable struct {
	errorAware
	pf ParamFuture
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
