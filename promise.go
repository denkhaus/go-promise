package promise

import (
	"fmt"
	"reflect"
)

type ParamFuture chan []reflect.Value
type ErrorFuture chan error

type Promise struct {
	pf ParamFuture
	ef ErrorFuture
}

///////////////////////////////////////////////////////////////////////////////////////
// makePromise
///////////////////////////////////////////////////////////////////////////////////////
func makePromise() *Promise {
	pr := &Promise{}
	pr.pf = make(ParamFuture)
	pr.ef = make(ErrorFuture)
	return pr
}

///////////////////////////////////////////////////////////////////////////////////////
// send
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) send(out []reflect.Value, err error) {
	p.pf <- out
	p.ef <- err
}

///////////////////////////////////////////////////////////////////////////////////////
// receive
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) receive() (out []reflect.Value, err error) {

	outR := false
	errR := false

	for !errR || !outR {
		select {
		case err = <-p.ef:
			errR = true
		case out = <-p.pf:
			outR = true
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////////////
// invoke
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) invoke(fn interface{}, in []reflect.Value) {
	v := reflect.ValueOf(fn)
	t := v.Type()

	//check arguments count equal
	if len(in) != t.NumIn() {
		// internal error, send the prev output and return internal error.
		p.send(in, fmt.Errorf("Function argument count mismatch."))
		return
	}
	//check arguments types equal
	for idx, inVal := range in {
		if inVal.Type() != t.In(idx) {
			// internal error, send the prev output and return internal error.
			p.send(in, fmt.Errorf("Function argument type mismatch."))
			return
		}
	}

	p.send(v.Call(in), nil)
}

///////////////////////////////////////////////////////////////////////////////////////
// Q
///////////////////////////////////////////////////////////////////////////////////////
func Q(init interface{}) *Promise {
	pr := makePromise()
	t := reflect.TypeOf(init)

	if t.Kind() == reflect.Func {
		go pr.invoke(init, []reflect.Value{})
	}

	return pr
}

///////////////////////////////////////////////////////////////////////////////////////
// Then
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Then(fns ...interface{}) *Promise {
	newP := makePromise()

	go func() {
		// wait on result and error from prev promise
		out, err := p.receive()

		// if we have an internal error, bubble it through, send the prev output and return.
		if err != nil {
			newP.send(out, err)
			return
		}

		// if we have only one Then func, map prev promise outputs to input
		if len(fns) == 1 {
			newP.invoke(fns[0], out)
			return
		}
	}()

	return newP
}

///////////////////////////////////////////////////////////////////////////////////////
// Done
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Done() ([]interface{}, error) {

	out, err := p.receive()
	res := make([]interface{}, len(out))

	for idx, val := range out {
		res[idx] = val.Interface()
	}

	return res, err
}
