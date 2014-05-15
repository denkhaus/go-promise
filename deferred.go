package Q

import (
	"reflect"
)

type Deferred struct {
	invokable
	targ []reflect.Value
	prev *Deferred
	next *Deferred
}

///////////////////////////////////////////////////////////////////////////////////////
// makeDeferred
///////////////////////////////////////////////////////////////////////////////////////
func makeDeferred(parent *Deferred, init []interface{}) *Deferred {
	df := &Deferred{prev: parent}
	df.targ = toValueArray(init)
	df.result = []reflect.Value{}
	df.rf = make(ResultFuture)

	if parent != nil {
		parent.next = df
	}

	return df
}

///////////////////////////////////////////////////////////////////////////////////////
// Defer creates a Deferred datatype. A Deferred can be resolved by value(s), func(s)
// Promise or Deferred.
///////////////////////////////////////////////////////////////////////////////////////
func Defer(init ...interface{}) *Deferred {
	df := makeDeferred(nil, init)
	return df
}

///////////////////////////////////////////////////////////////////////////////////////
// resolve
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) resolve() {

	in := []reflect.Value{}

	if d.prev != nil { // not the start element
		//receive new inputs from prev invocation
		in = d.prev.receive()
	}

	d.invokeTargets(d.targ, in)
}

///////////////////////////////////////////////////////////////////////////////////////
// Resolve
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Resolve(init ...interface{}) {

	first := func() *Deferred {
		start := d
		for start.prev != nil {
			start = start.prev
		}
		return start
	}

	dStart := first()
	// do we have inputs?
	if len(init) > 0 {
		// create start deferred to resolve inputs
		df := makeDeferred(nil, init)
		//link, and make it the very first
		df.next = dStart
		dStart.prev = df
		dStart = df
	}

	//resolve all deferreds
	for dStart != nil {
		go dStart.resolve()
		dStart = dStart.next
	}
}

///////////////////////////////////////////////////////////////////////////////////////
// Done
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Done() []interface{} {

	last := func() *Deferred {
		start := d
		for start.next != nil {
			start = start.next
		}

		return start
	}

	theLast := last()
	data := theLast.receive()
	return fromValueArray(data)
}

///////////////////////////////////////////////////////////////////////////////////////
// Reject
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Reject(err error) {

}

///////////////////////////////////////////////////////////////////////////////////////
// Deferred Then
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Then(init ...interface{}) *Deferred {
	df := makeDeferred(d, init)
	return df
}
