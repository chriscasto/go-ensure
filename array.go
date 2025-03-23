package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

// arrCheckFunc defines a function that can be used to validate an array
type arrCheckFunc[T any] func([]T) error

// ArrayValidator contains information and logic used to validate an array of type T
type ArrayValidator[T any] struct {
	typeStr      string
	lenValidator *NumberValidator[int]
	checks       []arrCheckFunc[T]
}

// Array constructs an ArrayValidator instance of type T and returns a pointer to it
func Array[T any]() *ArrayValidator[T] {
	var zero T

	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())

	return &ArrayValidator[T]{
		typeStr: typeStr,
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

// Each applies the provided validator to each element in an array and returns an error if any fail
func (v *ArrayValidator[T]) Each(ev with.Validator) *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		for _, e := range arr {
			if err := ev.Validate(e); err != nil {
				return err
			}
		}

		return nil
	})
}

// Validate accepts an arbitrary input type and validates it if it's a match for the expected type
func (v *ArrayValidator[T]) Validate(i interface{}) error {
	if err := testType(i, v.typeStr); err != nil {
		return err
	}
	return v.ValidateArray(i.([]T))
}

// ValidateArray applies all checks against an array and returns an error if any fail
func (v *ArrayValidator[T]) ValidateArray(arr []T) error {
	if v.lenValidator != nil {
		if err := v.lenValidator.Validate(len(arr)); err != nil {
			return err
		}
	}

	for _, fn := range v.checks {
		if err := fn(arr); err != nil {
			return err
		}
	}

	return nil
}

// Is adds the provided function as a check against any values to be validated
func (v *ArrayValidator[T]) Is(fn arrCheckFunc[T]) *ArrayValidator[T] {
	v.checks = append(v.checks, fn)
	return v
}
