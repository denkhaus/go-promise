package Q_test

import (
	"bitbucket.org/mendsley/tcgl/asserts"
	"fmt"
	"github.com/denkhaus/go-q"
	"testing"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestReturnValueIsNil
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestReturnValueIsEmptyAndErrorIsNil(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	Q.OnComposingError(func(e error) { err = e })

	res := Q.Promise(func() {}).Done()
	assert.Nil(err, "Error return value doesn't match.")
	assert.Empty(res, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestReturnValueIsValid1
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestReturnValueIsValid1(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	Q.OnComposingError(func(e error) { err = e })

	res := Q.Promise(func() int { return 5 }).Done()
	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestReturnValueIsValid2
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestReturnValueIsValid2(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	Q.OnComposingError(func(e error) { err = e })

	res := Q.Promise(func() (int, error) { return 5, fmt.Errorf("This is an error!") }).Done()
	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 2, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
	assert.ErrorMatch(res[1].(error), "This is an error!", "Error value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestBasicChainWithOneThenFunc
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestBasicChainWithOneThenFunc(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	var err error
	Q.OnComposingError(func(e error) { err = e })

	res := Q.Promise(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	}).Done()

	assert.Nil(err, "Error return value doesn't match.")
	assert.Length(res, 1, "Return value has invalid length.")
	assert.Equal(res[0], 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestBasicChainWithArgumentCountFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestBasicChainWithArgumentCountFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	Q.OnComposingError(func(err error) {
		assert.NotNil(err, "Error return value doesn't match.")
		assert.Substring(err.Error(), "Function argument count mismatch.", "Internal Error value doesn't match.")
	})

	res := Q.Promise(func() string {
		return "Hello Q"
	}).Then(func(theString string, theError error) int {

		//this will never be executed
		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	}).Done()

	assert.Length(res, 0, "Return value has invalid length.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestBasicChainWithArgumentTypeFailing
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestBasicChainWithArgumentTypeFailing(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	Q.OnComposingError(func(err error) {
		assert.NotNil(err, "Error return value doesn't match.")
		assert.Substring(err.Error(), "Function argument type mismatch.", "Internal Error value doesn't match.")
	})

	res := Q.Promise(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theError error, theString string) int {

		//this will never be executed
		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	}).Done()

	assert.Length(res, 0, "Return value has invalid length.")
}
