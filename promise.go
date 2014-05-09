package promise

import (
	//"fmt"
	"reflect"
)

type ParamFuture chan []reflect.Value
type ErrorFuture chan error

type Promise struct {
	pf ParamFuture
	ef ErrorFuture
}

func (p *Promise) dddd(fn interface{}) {
	v := reflect.ValueOf(fn)
	p.pf <- v.Call([]reflect.Value{})
	close(p.pf)

	p.ef <- nil
	close(p.ef)
}

///////////////////////////////////////////////////////////////////////////////////////
// Q
///////////////////////////////////////////////////////////////////////////////////////
func Q(fn interface{}) *Promise {
	pr := &Promise{}
	pr.pf = make(ParamFuture)
	pr.ef = make(ErrorFuture)

	go pr.dddd(fn)

	return pr
}

//func Defer(f func() interface{}) *Promise {
//	return &Promise{f: f}
//}

//func (p *Promise) Force() interface{} {
//	p.once.Do(func() { p.value = p.f() })
//	return p.value
//}

///////////////////////////////////////////////////////////////////////////////////////
// call
///////////////////////////////////////////////////////////////////////////////////////
//func (p *Promise) call() ([]reflect.Value, error) {

//	fn := p.fn
//fnType := reflect.TypeOf(fn)

//if len(in) != fnType.NumIn() {
//		return nil, fmt.Errorf("dfgjkl")
//}

//t.NumOut()
//t.In(i) // type of input variable i
//t.Out(i)

//	out := fn.Call(p.params)
//	return out, nil
//}
/*
///////////////////////////////////////////////////////////////////////////////////////
// Then
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Then(fns ...interface{}) *Promise {
	nPf := make(ParamFuture)
	nEf := make(ErrorFuture)

	go func(pf ParamFuture, ef ErrorFuture, fns []interface{}) {
		// wait on result and error from prev promise
		res := <-pf
		err := <-ef

		close(pf)
		close(ef)

		if err != nil {
			nEf <- nil
			nPf <- []reflect.Value{}
			return
		}

		// if we have only one Then func, map prev promise outputs to input
		if len(fns) == 1 {
			v := reflect.ValueOf(fns[0])
			if len(res) != v.Type().NumIn() {
				// outputs count from prev promise func mismatch current funcs inputs count
				nPf <- []reflect.Value{}
				nEf <- fmt.Errorf("Function argument count mismatch")
				return
			}

			nEf <- nil
			nPf <- v.Call(res)
			return
		}

	}(p.pf, p.ef, fns)

	newP := &Promise{pf: nPf, ef: nEf}
	return newP
}
*/
///////////////////////////////////////////////////////////////////////////////////////
// Do
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Done() (interface{}, error) {

	var (
		err error
		out interface{}
	)

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

	return out, err
}
