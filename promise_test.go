package promise

import (
	"bitbucket.org/mendsley/tcgl/asserts"
	"fmt"
	"testing"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestBasicChainWithError
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestBasicChainWithError(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	res := Q(func() (string, error) {
		return "Hello Q", fmt.Errorf("This is an error!")

	}).Then(func(theString string, theError error) int {

		assert.Equal(theString, "Hello Q", "String value doesn't match.")
		assert.ErrorMatch(theError, "This is an error!", "Error value doesn't match.")
		return 5

	}).Done()

	assert.Equal(res, 5, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestReturnValueIsNil
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestReturnValueIsNil(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	res := Q(func() {}).Done()
	assert.Nil(res, "Return value doesn't match.")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// TestReturnValueIsEqual
/////////////////////////////////////////////////////////////////////////////////////////////////////
func TestReturnValueIsEqual(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	res := Q(func() int { return 5 }).Done()
	assert.Equal(res, 5, "Return value doesn't match.")
}