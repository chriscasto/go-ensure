package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

// ArrayValidator contains information and logic used to validate an array of type T
type ArrayValidator[T any] struct {
	typeStr string
	checks  *iterChecks[int, T, []T]
}

// Array constructs an ArrayValidator instance of type T and returns a pointer to it
func Array[T any]() *ArrayValidator[T] {
	var zero T

	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())

	return &ArrayValidator[T]{
		typeStr: typeStr,
		checks:  newArrIterChecks[T](),
	}
}

// Type returns a string with the type of array this validator expects
func (v *ArrayValidator[T]) Type() string {
	return v.typeStr
}

// HasLengthWhere adds a NumberValidator for validating the length of the array
func (v *ArrayValidator[T]) HasLengthWhere(nv *NumberValidator[int]) *ArrayValidator[T] {
	v.checks.AddHasLengthWhere(nv)
	return v
}

// IsEmpty adds a check that returns an error if the length of the array is not 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (v *ArrayValidator[T]) IsEmpty() *ArrayValidator[T] {
	v.checks.AddIsEmpty()
	return v
}

// IsNotEmpty adds a check that returns an error if the length of the array is
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (v *ArrayValidator[T]) IsNotEmpty() *ArrayValidator[T] {
	v.checks.AddIsNotEmpty()
	return v
}

// HasCount adds a check that returns an error if the length of the array does not equal the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (v *ArrayValidator[T]) HasCount(l int) *ArrayValidator[T] {
	v.checks.AddHasLength(l)
	return v
}

// HasMoreThan adds a check that returns an error if the length of the array is less than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsGreaterThan(l))
func (v *ArrayValidator[T]) HasMoreThan(l int) *ArrayValidator[T] {
	v.checks.AddIsLongerThan(l)
	return v
}

// HasFewerThan adds a check that returns an error if the length of the array is more than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsLessThan(l))
func (v *ArrayValidator[T]) HasFewerThan(l int) *ArrayValidator[T] {
	v.checks.AddIsShorterThan(l)
	return v
}

// Each assigns a Validator to be used for validating array values
func (v *ArrayValidator[T]) Each(ev with.Validator[T]) *ArrayValidator[T] {
	v.checks.AddIterValValidator(ev)
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
	return v.checks.Evaluate(arr, getValidationOptions(options))
}

// Is adds the provided function as a check against any values to be validated
func (v *ArrayValidator[T]) Is(fn func([]T) error) *ArrayValidator[T] {
	v.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		return fn(val)
	})
	return v
}

// Has adds the provided function as a check against any values to be validated
// Has is an alias for Is
func (v *ArrayValidator[T]) Has(fn func([]T) error) *ArrayValidator[T] {
	return v.Is(fn)
}
