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
func DeferredTestReturnValueIsEmptyAndErrorIsNil(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func() {})
	df.Resolve()
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Empty(res, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestReturnValueIsValid1
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestReturnValueIsValid1(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func() int { return 5 })
	df.Resolve()
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestReturnValueIsValid2
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestReturnValueIsValid2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func() (int, error) { return 5, fmt.Errorf("This is an error!") })
	df.Resolve()
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
	assert.ErrorMatch(res[1].(error), "This is an error!", "Error value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithOneThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithOneThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	})

	df.Resolve()
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithArgumentCountFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithArgumentCountFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

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
	res, err := df.Done()

	assert.False(executed, "Then func was executed.")
	assert.NotNil(err, "Error return value doesn't match.")
	assert.Substring(err.Error(), "Function argument count mismatch.", "Internal Error value doesn't match.")
	assert.Length(res, 0, "Return value has invalid length.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithArgumentTypeFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithArgumentTypeFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

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
	res, err := df.Done()

	assert.False(executed, "Then func was executed.")
	assert.NotNil(err, "Error return value doesn't match.")
	assert.Substring(err.Error(), "Function argument type mismatch.", "Internal Error value doesn't match.")
	assert.Length(res, 0, "Return value has invalid length.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithFirstInputFromResolve1
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithFirstInputFromResolve1(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5
	})

	df.Resolve("Hello Q")
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithFirstInputFuncFromResolve1
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithFirstInputFuncFromResolve1(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func(inp string) (string, error) {
		return inp, fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5
	})

	df.Resolve(func() string { return "Hello Q" })
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithFirstInputFromResolve
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithFirstInputFromResolve2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func(inp1 string, inp2 int) (string, int) {
		return inp1, inp2

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5
	})

	df.Resolve("Hello Q", 10)
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithFirstInputFuncFromResolve2
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithFirstInputFuncFromResolve2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	df := Defer(func(inp1 string, inp2 int) (string, int) {
		return inp1, inp2

	}).Then(func(theString string, theInt int) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.Equal(theInt, 10, "Int value doesn't match.")
		return 5
	})

	df.Resolve(func() (string, int) { return "Hello Q", 10 })
	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// DeferredTestBasicChainWithFirstInputFuncFromResolve2ThenFuncs
/////////////////////////////////////////////////////////////////////////////////////////////////////
func DeferredTestBasicChainWithFirstInputFuncFromResolve2ThenFuncs(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

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

	res, err := df.Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value1 doesn't match.")
	assert.Equal(res[0], 10, "Return value2 doesn't match.")
}
