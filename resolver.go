package Q

import (
	"errors"
	"fmt"
	"reflect"
)

type resolver struct {
	pr  Progressor
	val reflect.Value
	t   reflect.Type
}

type OnResolveFunc func([]reflect.Value)

///////////////////////////////////////////////////////////////////////////////////////
// CanInvokeWithParams checks if func is invokable by given input arguments.
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) CanInvokeWithParams(in []reflect.Value) bool {

	argCount := r.InArgCount()
	if len(in) < argCount {
		return false
	}

	for i := 0; i < argCount; i++ {
		if !in[i].Type().AssignableTo(r.InArgType(i)) {
			return false
		}
	}

	return true
}

///////////////////////////////////////////////////////////////////////////////////////
// Resolver initializes and returns a new Resolver object
///////////////////////////////////////////////////////////////////////////////////////
func Resolver(i *invokable, val reflect.Value) *resolver {
	r := &resolver{val: val, pr: i.pr, t: val.Type()}
	return r
}

///////////////////////////////////////////////////////////////////////////////////////
// Resolve() composes input variables, injects Q.Progressor if required, resolves Q.Deferred
// an Q.Promised inputs and maps result to func arguments.
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) Resolve(in []reflect.Value, onResolve OnResolveFunc) ([]reflect.Value, error) {

	var (
		inIdx int
		remIn = in
		nOut  = r.InArgCount()
	)

	resOut := []reflect.Value{}
	for outIdx := 0; outIdx < nOut; outIdx++ {
		outType := r.InArgType(outIdx)
		if reflect.TypeOf(r.pr).AssignableTo(outType) {
			resOut = append(resOut, reflect.ValueOf(r.pr))
			continue
		}

		if inIdx >= len(in) {
			break
		}

		inType := in[inIdx].Type()
		v, ok := in[inIdx].Interface().(Invokable)
		if !inType.AssignableTo(outType) && ok {
			res := v.receive()
			resOut = append(resOut, res...)
			outIdx += len(res)
		} else {
			resOut = append(resOut, in[inIdx])
		}

		inIdx++
		remIn = in[inIdx:]
	}

	if r.CanInvokeWithParams(resOut) {
		onResolve(resOut)
		return remIn, nil
	} else {

		//do we have enough func inputs?
		if len(resOut) < nOut {
			return nil, errors.New("Function argument count mismatch. Need more inputs.")
		}

		//check for argument errors
		for idx, inVal := range resOut {
			t := r.InArgType(idx)
			if inVal.Type() != t {
				return nil, fmt.Errorf("Function argument type mismatch. (%v -> %v)", inVal.Type(), t)
			}
		}
	}

	return nil, nil
}

///////////////////////////////////////////////////////////////////////////////////////
// InArgCount
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) InArgCount() int {
	return r.t.NumIn()
}

///////////////////////////////////////////////////////////////////////////////////////
// InArgType
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) InArgType(n int) reflect.Type {
	return r.t.In(n)
}
