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
	for actIdx, actInp := range in {
		actInpType := actInp.Type()
		targIdx := actIdx + delta

		if targIdx >= nFnInpts {
			break
		}

		v, ok := actInp.Interface().(Invokable)
		if actInpType != r.InArgType(targIdx) && ok {
			res := v.receive()
			resInp = append(resInp, res...)
			delta += len(res)
		} else {
			resInp = append(resInp, actInp)
		}
		remInp = in[actIdx+1:]
	}

	//check again
	if r.CanInvokeWithParams(resInp) {
		onResolve(resInp)
		return remInp, nil
	} else {

		//check we have enough func inputs
		if len(resInp) < nFnInpts {
			return nil, errors.New("Function argument count mismatch. Need more inputs.")
		}

		//check for argument errors
		for idx, inVal := range resInp {
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
