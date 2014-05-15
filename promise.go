package Q

import (
	"reflect"
)

type Promised struct {
	invokable
}

///////////////////////////////////////////////////////////////////////////////////////
// makePromise
///////////////////////////////////////////////////////////////////////////////////////
func makePromised() *Promised {
	pr := new(Promised)
	pr.rf = make(ResultFuture)
	pr.result = []reflect.Value{}

	return pr
}

///////////////////////////////////////////////////////////////////////////////////////
// Promise Then
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Then(init ...interface{}) *Promised {
	newP := makePromised()

	go func() {
		// old error from prev promises
		if p.err != nil {
			newP.sendError(nil, 0, p.err)
			return
		}

		// wait on result from prev promise
		in := p.receive()
		// and invoke it
		targets := toValueArray(init)
		newP.invokeTargets(targets, in)
	}()

	return newP
}

///////////////////////////////////////////////////////////////////////////////////////
// Done
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Done() []interface{} {

	out := p.receive()
	res := fromValueArray(out)
	return res
}

///////////////////////////////////////////////////////////////////////////////////////
// Finally
///////////////////////////////////////////////////////////////////////////////////////
func (p *Promised) Finally(init interface{}) {

	if p.err != nil {
		return
	}

	in := p.receive()
	//TODO change that, use other toValueArray version
	vals := make([]interface{}, 1)
	vals[0] = init

	targets := toValueArray(vals)
	go p.invokeTargets(targets, in)
}

///////////////////////////////////////////////////////////////////////////////////////
// Promise
///////////////////////////////////////////////////////////////////////////////////////
func Promise(init ...interface{}) *Promised {

	pr := makePromised()
	targets := toValueArray(init)
	go pr.invokeTargets(targets, nil)
	return pr
}
