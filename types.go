package ensure

import (
	"fmt"
	"reflect"
)

// Length is a convenience function for creating a validator to be used on types with a length property
func Length() *NumberValidator[int] {
	return Number[int]()
}

// TypeError indicates a mismatch between the type expected by a validator and
// the type of the value passed to the validator.  This should generally not
// be passed back to the user.
type TypeError struct {
	err string
}

func (e *TypeError) Error() string {
	return e.err
}

func NewTypeError(err string) *TypeError {
	return &TypeError{err}
}

// newTypeErrorFromTypes is a helper function for creating a TypeError in the
// common scenario where you have the names of the type you want and the type
// you have already available in string form
func newTypeErrorFromTypes(want string, got string) *TypeError {
	return NewTypeError(fmt.Sprintf(
		`expected "%s"; got "%s"`, want, got,
	))
}

// ValidationError is used to indicate a failure while conducting validation
// checks.  These are generally safe to return to the user so they can correct
// their input(s)
type ValidationError struct {
	err string
}

// NewValidationError returns a ValidationError with the error message passed to it
func NewValidationError(err string) *ValidationError {
	return &ValidationError{err}
}

func (e *ValidationError) Error() string {
	return e.err
}

// testType compares a value against an expected type and returns a type error if they don't match
func testType(value any, expect string) *TypeError {
	valType := reflect.TypeOf(value).String()

	if valType != expect {
		return newTypeErrorFromTypes(expect, valType)
	}

	return nil
}
