package promise

import (
	//"fmt"
	"reflect"
)

type Param reflect.Value
type Params []Param

type Promise struct {
	f chan Params
}

///////////////////////////////////////////////////////////////////////////////////////
// Q
///////////////////////////////////////////////////////////////////////////////////////
func Q(fn interface{}) *Promise {
	future := make(chan Params)

	go func() {
		v := reflect.ValueOf(fn)
		future <- v.Call()
	}()

	return &Promise{f: future}
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
func (p *Promise) call() ([]reflect.Value, error) {

	fn := p.fn
	//fnType := reflect.TypeOf(fn)

	//if len(in) != fnType.NumIn() {
	//		return nil, fmt.Errorf("dfgjkl")
	//}

	//t.NumOut()
	//t.In(i) // type of input variable i
	//t.Out(i)

	out := fn.Call(p.params)
	return out, nil
}

///////////////////////////////////////////////////////////////////////////////////////
// Then
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Then(fns ...interface{}) *Promise {

	if len(fns) == 1 {

	}

	newP := &Promise{fn: fn, err: p.err}

	if p.err != nil {
		return newP
	}

	out, err := p.call()
	newP.params = out
	newP.err = err

	return newP
}

///////////////////////////////////////////////////////////////////////////////////////
// Do
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promise) Do() (interface{}, error) {

	if p.err != nil {
		return nil, p.err
	}

	out, err := p.call()
	return out, err
}
