package ensure

import (
	"fmt"
	"reflect"
)

// Length is a convenience function for creating a validator to be used on types with a length property
func Length() *NumberValidator[int] {
	return Number[int]()
}

// newTypeErrorFromTypes is a helper function for creating a TypeError in the
// common scenario where you have the names of the type you want and the type
// you have already available in string form
func newTypeErrorFromTypes(want string, got string) *TypeError {
	return NewTypeError(fmt.Sprintf(
		`expected "%s"; got "%s"`, want, got,
	))
}

// testType compares a value against an expected type and returns a type error if they don't match
func testType(value any, expect string) *TypeError {
	valType := reflect.TypeOf(value).String()

	if valType != expect {
		return newTypeErrorFromTypes(expect, valType)
	}

	return nil
}
