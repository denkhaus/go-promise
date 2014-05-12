package Q

import (
	"reflect"
)

type Promised struct {
	invokable
}

type Deferred struct {
	invokable
	inputs []reflect.Value
	prev   *Deferred
	next   *Deferred
}

///////////////////////////////////////////////////////////////////////////////////////
// makeDeferred
///////////////////////////////////////////////////////////////////////////////////////
func makeDeferred(parent *Deferred) *Deferred {
	df := &Deferred{prev: parent}
	if parent != nil {
		parent.next = df
	}

	return df
}

///////////////////////////////////////////////////////////////////////////////////////
// resolve
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) resolve(inp []interface{}) {

	var params []reflect.Value

	if d.inputs == nil {
		params = []reflect.Value{}

		if inp != nil {
			//set start inputs
			d.inputs = toValueArray(inp)
		}
	} else {
		params = d.receive()
	}

	in := d.inputs

	if len(in) == 1 {
		t := in[0].Type()

		//TODO handle promise input

		if t.Kind() == reflect.Func {
			//input is func, invoke it
			go d.invoke(in[0], params)
		} else {
			//TODO need testcase
			//input is init value, send it directly
			go d.send(in)
		}
	} else {
		//TODO need testcase
		//inputs are init values, send it directly
		go d.send(in)
	}
}

///////////////////////////////////////////////////////////////////////////////////////
// makePromise
///////////////////////////////////////////////////////////////////////////////////////
func makePromised() *Promised {
	pr := new(Promised)
	pr.pf = make(ParamFuture)
	return pr
}

///////////////////////////////////////////////////////////////////////////////////////
// Defer creates a Deferred datatype. A Deferred can be resolved by value or a promise.
///////////////////////////////////////////////////////////////////////////////////////
func Defer() *Deferred {
	df := makeDeferred(nil)
	return df
}

///////////////////////////////////////////////////////////////////////////////////////
// Resolve |
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Resolve(init ...interface{}) {
	//TODO handle promise input

	dStart := d
	for dStart.prev != nil {
		dStart = dStart.prev
	}
	//start at the beginning of chain
	go dStart.resolve(init)
}

///////////////////////////////////////////////////////////////////////////////////////
// Reject
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Reject(err error) {

}

///////////////////////////////////////////////////////////////////////////////////////
// Deferred Then
///////////////////////////////////////////////////////////////////////////////////////
func (d *Deferred) Then(fns ...interface{}) *Deferred {
	df := makeDeferred(d)
	df.inputs = toValueArray(fns)
	return df
}

///////////////////////////////////////////////////////////////////////////////////////
// Promise Then
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Then(fns ...interface{}) *Promised {
	newP := makePromised()

	go func() {
		// old error from prev promises
		if p.err != nil {
			newP.send([]reflect.Value{}) // send dummy to avoid goroutine deadlock
			newP.setError(nil, p.err.Error())
			return
		}

		// wait on result from prev promise
		out := p.receive()

		// if we have only one Then func, map prev promise outputs to input
		if len(fns) == 1 {
			//TODO handle promise input
			newP.invoke(fns[0], out)
		}
	}()

	return newP
}

///////////////////////////////////////////////////////////////////////////////////////
// Done
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Done() ([]interface{}, error) {

	out := p.receive()
	res := fromValueArray(out)
	return res, p.err
}

///////////////////////////////////////////////////////////////////////////////////////
// Finally
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Finally(fn interface{}) error {

	if p.err != nil {
		return p.err
	}

	out := p.receive()
	t := reflect.TypeOf(fn)
	if t.Kind() == reflect.Func {
		go p.invoke(fn, out)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////
// Promise
///////////////////////////////////////////////////////////////////////////////////////
func Promise(init ...interface{}) *Promised {

	pr := makePromised()

	if len(init) == 1 {
		t := reflect.TypeOf(init[0])

		//TODO handle promise input

		if t.Kind() == reflect.Func {
			//input is init func, invoke it
			go pr.invoke(init[0], []reflect.Value{})
		} else {
			//TODO need testcase
			//input is init value, send it directly
			v := reflect.ValueOf(init[0])
			go pr.send([]reflect.Value{v})
		}
	} else {
		//TODO need testcase
		//inputs are init values, send it directly
		out := toValueArray(init)
		go pr.send(out)
	}

	return pr
}
