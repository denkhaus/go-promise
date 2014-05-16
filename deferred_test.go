package Q

import (
	"bitbucket.org/mendsley/tcgl/asserts"
	"errors"
	"fmt"
	"testing"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestReturnValueIsNil
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredReturnValueIsEmptyAndErrorIsNil(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	executed := 0
	df := Defer(func() { executed++ })
	df.Resolve()

	res := df.Done()

	assert.Equal(executed, 1, "Executed value doesn't match.")
	assert.Nil(err, "Error return value doesn't match.")
	assert.Empty(res, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredReturnValueIsValid1
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredReturnValueIsValid1(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func() int { return 5 })
	df.Resolve()
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestReturnValueIsValid2
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredReturnValueIsValid2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func() (int, error) { return 5, fmt.Errorf("This is an error!") })
	df.Resolve()
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
	assert.ErrorMatch(res[1].(error), "This is an error!", "Error value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithOneThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithOneThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	})

	df.Resolve()
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainResolveByValue
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainResolveByValue(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")
	})

	df.Resolve("Hello Q")
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], "Hello Q", "Return value doesn't match.")
	assert.ErrorMatch(res[1].(error), "This is an error!", "Error value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainResolveByFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainResolveByFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")
	})

	df.Resolve(func() string { return "Hello Q" })
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], "Hello Q", "Return value doesn't match.")
	assert.ErrorMatch(res[1].(error), "This is an error!", "Error value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithArgumentCountFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithArgumentCountFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	OnInternalError(func(err error) {
		assert.NotNil(err, "Error return value doesn't match.")
		assert.Substring(err.Error(), "Function argument count mismatch.", "Internal Error value doesn't match.")
	})

	executed := false
	df := Defer(func() string {
		return "Hello Q"
	}).Then(func(theString string, theError error) int {
		executed = true
		//this will never be executed
		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	})

	df.Resolve()
	res := df.Done()

	assert.False(executed, "Then func was executed.")
	assert.Length(res, 0, "Return value has invalid length.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithArgumentTypeFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithArgumentTypeFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	OnInternalError(func(err error) {
		assert.NotNil(err, "Error return value doesn't match.")
		assert.Substring(err.Error(), "Function argument type mismatch.", "Internal Error value doesn't match.")
	})

	executed := false
	df := Defer(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theError error, theString string) int {
		executed = true
		//this will never be executed
		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	})

	df.Resolve()
	res := df.Done()

	assert.False(executed, "Then func was executed.")
	assert.Length(res, 0, "Return value has invalid length.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithUnusedInputsErrorReporting
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithUnusedInputsErrorReporting(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	OnInternalError(func(err error) {
		assert.NotNil(err, "Error return value doesn't match.")
		assert.Substring(err.Error(), "Unused inputs on target.", "Internal Error value doesn't match.")
	})

	executed := false
	df := Defer(func() (string, int, error) {
		return "Hello Q", 10, fmt.Errorf("This is an error!")

	}).Then(func(theString string, theInt int) int {
		executed = true
		//this will never be executed
		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Integer value doesn't match.")

		return 5
	})

	df.Resolve()
	res := df.Done()

	assert.True(executed, "Then func was executed.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithResolveByValueAndThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithResolveByValueAndThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5
	})

	df.Resolve("Hello Q")
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithResolveByFuncAndThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithResolveByFuncAndThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5
	})

	df.Resolve(func() string { return "Hello Q" })
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithResolveByValuesAndThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithResolveByValuesAndThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp1 string, inp2 int) (string, int) {
		return inp1, inp2

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5
	})

	df.Resolve("Hello Q", 10)
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithResolveByPromiseAndThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithResolveByPromisAndThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp1 string, inp2 int) (string, int) {
		return inp1, inp2

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5
	})

	p := Promise(func() (string, int) {
		return "Hello Q", 10
	})

	df.Resolve(p)
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithResolveByFuncAndThenFunc2
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithResolveByFuncAndThenFunc2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp1 string, inp2 int) (string, int) {
		return inp1, inp2

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5
	})

	df.Resolve(func() (string, int) { return "Hello Q", 10 })
	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithOverlapping
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithOverlapping(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp1 string, inp2 int, errorText string) (string, int, error) {
		return inp1, inp2, errors.New(errorText)

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5

	}, func(err error) int {
		assert.ErrorMatch(err, "This is an error!", "Error value doesn't match.")
		return 10
	})

	df.Resolve(func() (string, int, string) {
		return "Hello Q", 10, "This is an error!"
	})

	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value1 doesn't match.")
	assert.Equal(res[1], 10, "Return value2 doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestDeferredBasicChainWithOverlappingAndPromiseInput
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDeferredBasicChainWithOverlappingAndPromiseInput(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	OnInternalError(func(e error) { err = e })

	df := Defer(func(inp1 string, inp2 int, errorText string) (string, *Promised, error) {

		p := Promise(func() int {
			return inp2
		})

		return inp1, p, errors.New(errorText)

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5

	}, func(err error) int {
		assert.ErrorMatch(err, "This is an error!", "Error value doesn't match.")
		return 10
	})

	df.Resolve(func() (string, int, string) {
		return "Hello Q", 10, "This is an error!"
	})

	res := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value1 doesn't match.")
	assert.Equal(res[1], 10, "Return value2 doesn't match.")
}
