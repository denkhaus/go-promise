package Q

import (
	"errors"
	"fmt"
	"reflect"
)

type resolver struct {
	pr  reflect.Value
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
func Resolver(i *invokable, val reflect.Value) *resolver {
	r := &resolver{val: val,
		pr: reflect.ValueOf(i.pr),
		t:  val.Type()}
	return r
}

///////////////////////////////////////////////////////////////////////////////////////
// Resolve
///////////////////////////////////////////////////////////////////////////////////////
func (r *resolver) Resolve(in []reflect.Value, onResolve OnResolveFunc) ([]reflect.Value, error) {

	inIdx := 0
	remIn := in
	nOut := r.InArgCount()
	f := r.pr.Interface()
	resOut := []reflect.Value{}

	for outIdx := 0; outIdx < nOut; outIdx++ {
		outType := r.InArgType(outIdx)
		if reflect.TypeOf(f).AssignableTo(outType) {
			resOut = append(resOut, r.pr)
			continue
		}

		if inIdx >= nOut {
			break
		}

		inType := in[inIdx].Type()
		if inType != outType {
			if v, ok := in[inIdx].Interface().(Invokable); ok {
				res := v.receive()
				resOut = append(resOut, res...)
				outIdx += len(res)
			}
		} else {
			resOut = append(resOut, in[inIdx])
		}
		remInp = in[inIdx+1:]
		inIdx++
	}

	//if len(in) < nFnInpts {

	//	f := r.pr.Interface()

	//	for i := 0; i < nFnInpts; i++ {
	//		if reflect.TypeOf(f).AssignableTo(r.InArgType(i)) {
	//			fmt.Print("ffffffffffffff")
	//		}
	//	}
	//}

	//delta := 0
	//resInp := []reflect.Value{}
	//remInp := []reflect.Value{}
	//for actIdx, actInp := range in {
	//	actInpType := actInp.Type()
	//	targIdx := actIdx + delta

	//	if targIdx >= nFnInpts {
	//		break
	//	}

	//	actOutType := r.InArgType(targIdx)
	//	resolved := false
	//	if actInpType != actOutType {
	//		fmt.Print(actOutType)
	//		if v, ok := actInp.Interface().(Invokable); ok {
	//			res := v.receive()
	//			resInp = append(resInp, res...)
	//			delta += len(res)
	//			resolved = true
	//		} else {
	//			//fmt.Print(actOutType)
	//			// if actOutType.Name() == "Q.Progressor" {
	//			//resInp = append(resInp, r.pr)
	//			fmt.Print("ffffffffffffffffffffffffffffffffffffffff")
	//			//resolved = true
	//			//delta++
	//		}

	//		fmt.Print(actOutType)

	//	}

	//	if !resolved {
	//		resInp = append(resInp, actInp)
	//	}

	//	remInp = in[actIdx+1:]

	//if actInpType != actOutType {
	//	if v, ok := actInp.Interface().(Invokable); ok {
	//		res := v.receive()
	//		resInp = append(resInp, res...)
	//		delta += len(res)

	//	} /*else if actOutType == reflect.TypeOf(r.pr.Interface()) {
	//		resInp = append(resInp, r.pr)
	//		resInp = append(resInp, actInp)
	//		print("ffffffffffffffffffffffffffffffffffffffff")
	//		delta++
	//	}*/

	//} else {
	//	resInp = append(resInp, actInp)
	//}

	//}

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
