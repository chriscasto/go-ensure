package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
	"slices"
)

// ArrayValidator contains information and logic used to validate an array of type T
type ArrayValidator[T any] struct {
	typeStr       string
	lenValidator  *NumberValidator[int]
	itemValidator with.Validator[T]
	//itemChecks    *seqChecks[T]
	checks *valChecks[[]T]
}

// Array constructs an ArrayValidator instance of type T and returns a pointer to it
func Array[T any]() *ArrayValidator[T] {
	var zero T

	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())

	return &ArrayValidator[T]{
		typeStr: typeStr,
		checks:  newValidationChecks[[]T](),
	}
}

// Type returns a string with the type of array this validator expects
func (v *ArrayValidator[T]) Type() string {
	return v.typeStr
}

// HasLengthWhere adds a NumberValidator for validating the length of the array
func (v *ArrayValidator[T]) HasLengthWhere(nv *NumberValidator[int]) *ArrayValidator[T] {
	v.lenValidator = nv
	return v
}

// IsEmpty adds a check that returns an error if the length of the array is not 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (v *ArrayValidator[T]) IsEmpty() *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) != 0 {
			return errors.New(`array must be empty`)
		}

		return nil
	})
}

// IsNotEmpty adds a check that returns an error if the length of the array is
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (v *ArrayValidator[T]) IsNotEmpty() *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) == 0 {
			return errors.New(`array must not be empty`)
		}

		return nil
	})
}

// HasCount adds a check that returns an error if the length of the array does not equal the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (v *ArrayValidator[T]) HasCount(l int) *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) != l {
			return errors.New(
				fmt.Sprintf(
					`array length must equal %d; got %d`,
					l,
					len(arr)),
			)
		}

		return nil
	})
}

// HasMoreThan adds a check that returns an error if the length of the array is less than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsGreaterThan(l))
func (v *ArrayValidator[T]) HasMoreThan(l int) *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) <= l {
			return errors.New(
				fmt.Sprintf(
					`array must have a length greater than %d; got %d`,
					l,
					len(arr)),
			)
		}

		return nil
	})
}

// HasFewerThan adds a check that returns an error if the length of the array is more than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsLessThan(l))
func (v *ArrayValidator[T]) HasFewerThan(l int) *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) >= l {
			return errors.New(
				fmt.Sprintf(
					`array must have a length less than %d; got %d`,
					l,
					len(arr)),
			)
		}

		return nil
	})
}

// Each assigns a Validator to be used for validating array values
func (v *ArrayValidator[T]) Each(ev with.Validator[T]) *ArrayValidator[T] {
	v.itemValidator = ev
	return v
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a match for the expected type
func (v *ArrayValidator[T]) ValidateUntyped(value any, options ...*with.ValidationOptions) error {
	if err := testType(value, v.typeStr); err != nil {
		return err
	}
	return v.Validate(value.([]T), options...)
}

// Validate applies all checks against an array and returns an error if any fail
func (v *ArrayValidator[T]) Validate(arr []T, options ...*with.ValidationOptions) error {
	vOpts := getValidationOptions(options)
	abortOnErr := !vOpts.CollectAllErrors()

	var lenCheck func([]T) error

	// if there is a length validator, add another check to validate length with current opts
	if v.lenValidator != nil {
		lenCheck = func(arrVal []T) error {
			if err := v.lenValidator.Validate(len(arrVal), vOpts); err != nil {
				return err
			}
			return nil
		}
	}

	var eachCheck func([]T) error

	// if we have an item validator, add a check for each item in the array
	if v.itemValidator != nil {
		eachCheck = func(arrVal []T) error {
			itemSeqChecks := newSeqChecks[T](func(val T) error {
				if err := v.itemValidator.Validate(val, vOpts); err != nil {
					return err
				}

				return nil
			})

			seq := slices.Values(arrVal)

			if abortOnErr {
				return itemSeqChecks.Evaluate(seq)
			}

			return itemSeqChecks.EvaluateAll(seq)
		}
	}

	if vOpts.CollectAllErrors() {
		if err := v.checks.EvaluateAll(arr, lenCheck, eachCheck); err != nil {
			if err.HasErrors() {
				return err
			}
		}
	} else {
		return v.checks.Evaluate(arr, lenCheck, eachCheck)
	}

	return nil
}

// Is adds the provided function as a check against any values to be validated
func (v *ArrayValidator[T]) Is(fn func([]T) error) *ArrayValidator[T] {
	v.checks.Append(fn)
	return v
}
