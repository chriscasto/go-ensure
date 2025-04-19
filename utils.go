package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
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

func getValidationOptions(options []*with.ValidationOptions) *with.ValidationOptions {
	if len(options) > 0 {
		return options[0]
	}

	// default options
	return with.Options()
}

type validationChecks[T any] struct {
	c []func(T) error
}

func newValidationChecks[T any]() *validationChecks[T] {
	return &validationChecks[T]{}
}

func (vc *validationChecks[T]) Append(check func(T) error) {
	vc.c = append(vc.c, check)
}

func (vc *validationChecks[T]) Evaluate(target T, abortOnErr bool) error {
	// The code could be simplified by putting the abortOnErr check in the loop,
	// but since check evaluation is abstracted away anyway and the loop logic
	// is simple, I'm fine with a little over-optimization

	if abortOnErr {
		for _, fn := range vc.c {
			if err := fn(target); err != nil {
				return err
			}
		}
	} else {
		vErrs := newValidationErrors()

		for _, fn := range vc.c {
			if err := fn(target); err != nil {
				vErrs.Append(err)
			}
		}

		if vErrs.HasErrors() {
			return vErrs
		}
	}

	return nil
}
