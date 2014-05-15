package Q

import (
	"errors"
	"fmt"
	"reflect"
)

type resolver struct {
	val reflect.Value
	t   reflect.Type
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

///////////////////////////////////////////////////////////////////////////////////////
// CanInvokeWithParams
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) CanInvokeWithParams(in []reflect.Value) bool {

	argCount := r.InArgCount()
	if len(in) < argCount {
		return false
	}

	for i := 0; i < argCount; i++ {
		if in[i].Type() != r.InArgType(i) {
			return false
		}
	}

	return true
}

///////////////////////////////////////////////////////////////////////////////////////
// IsResolvable
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) IsResolvable(t reflect.Type) bool {

	if t == DeferredPtrType || t == PromisedPtrType {

		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////
// InvokeHelper
///////////////////////////////////////////////////////////////////////////////////////
func Resolver(val reflect.Value) *resolver {
	r := &resolver{val: val}
	r.t = val.Type()
	return r
}

type OnResolveFunc func([]reflect.Value)

///////////////////////////////////////////////////////////////////////////////////////
// Resolve
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) Resolve(in []reflect.Value, onResolve OnResolveFunc) ([]reflect.Value, error) {

	nFnInpts := r.InArgCount()

	//check we have enough func inputs
	if len(in) < nFnInpts {
		return nil, errors.New("Function argument count mismatch. Need more inputs.")
	}

	inIdx := 0
	resInput := []reflect.Value{}
	for resIdx := 0; resIdx < nFnInpts; resIdx++ {
		inpType := in[inIdx].Type()

		if inpType != r.InArgType(resIdx) && r.IsResolvable(inpType) {

			v := in[inIdx].Interface().(*Promised)
			res := v.receive()
			resInput = append(resInput, res...)
			resIdx += len(res)
		} else {
			resInput = append(resInput, in[inIdx])
		}
		inIdx++
	}

	//check again
	if r.CanInvokeWithParams(resInput) {
		onResolve(resInput)
		return in[inIdx:], nil
	} else {

		//check for argument errors
		for idx, inVal := range resInput {
			t := r.InArgType(idx)
			if inVal.Type() != t {
				return nil, fmt.Errorf("Function argument type mismatch. (%v -> %v)", inVal.Type(), t)
			}
		}
	}

	return nil, nil
}
