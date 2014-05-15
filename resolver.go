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
type OnResolveFunc func([]reflect.Value)

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

///////////////////////////////////////////////////////////////////////////////////////
// Resolve
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) Resolve(in []reflect.Value, onResolve OnResolveFunc) ([]reflect.Value, error) {

	nFnInpts := r.InArgCount()

	delta := 0
	resInp := []reflect.Value{}
	remInp := []reflect.Value{}
	for actIdx, actInput := range in {
		actInpType := actInp.Type()
		targIdx := actIdx + delta

		if targIdx >= nFnInpts {
			break
		}

		if actInputType != r.InArgType(targIdx) &&
			r.IsResolvable(actInputType) {
			v := actInput.Interface().(*Promised)
			res := v.receive()
			resInput = append(resInput, res...)
			delta += len(res)
		} else {
			resInput = append(resInput, actInput)
		}
		remInp = in[actIdx+1:]
	}

	//check again
	if r.CanInvokeWithParams(resInput) {
		onResolve(resInput)
		return remInput, nil
	} else {

		//check we have enough func inputs
		if len(resInput) < nFnInpts {
			return nil, errors.New("Function argument count mismatch. Need more inputs.")
		}

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
