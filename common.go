package Q

import (
	"reflect"
)

var (
	DeferredPtrType = reflect.TypeOf(&Deferred{})
	PromisedPtrType = reflect.TypeOf(&Promised{})
)

///////////////////////////////////////////////////////////////////////////////////////
// toValueArray
///////////////////////////////////////////////////////////////////////////////////////
func fromValueArray(in []reflect.Value) []interface{} {

	out := make([]interface{}, len(in))

	for idx, val := range in {
		out[idx] = val.Interface()
	}

	return out
}

///////////////////////////////////////////////////////////////////////////////////////
// toValueArray
///////////////////////////////////////////////////////////////////////////////////////

func toValueArray(in []interface{}) []reflect.Value {
	out := make([]reflect.Value, len(in))
	for idx, val := range in {
		out[idx] = reflect.ValueOf(val)
	}

	return out
}

/*
whats wrong with that?
func toValueArray(in ...interface{}) []reflect.Value {

	out := []reflect.Value{}

	for _, val := range in {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				out = append(out, v.Index(i))
			}
		} else {
			out = append(out, v)
		}
	}

	return out
}
*/
