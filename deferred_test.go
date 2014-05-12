package Q

import (
	"bitbucket.org/mendsley/tcgl/asserts"
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
