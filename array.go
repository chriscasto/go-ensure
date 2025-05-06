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

// ComparableArrayValidator validates arrays of comparable types
type ComparableArrayValidator[T comparable] struct {
	ArrayValidator[T]
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

// ComparableArray constructs a ComparableArrayValidator instance of type T and returns a pointer to it
func ComparableArray[T comparable]() *ComparableArrayValidator[T] {
	var zero T

	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())

	return &ComparableArrayValidator[T]{
		ArrayValidator[T]{
			typeStr: typeStr,
			checks:  newArrIterChecks[T](),
		},
	}
}

// Type returns a string with the type of array this validator expects
func (av *ArrayValidator[T]) Type() string {
	return av.typeStr
}

// HasLengthWhere adds a NumberValidator for validating the length of the array
func (av *ArrayValidator[T]) HasLengthWhere(nv *NumberValidator[int]) *ArrayValidator[T] {
	av.checks.AddHasLengthWhere(nv)
	return av
}

// IsEmpty adds a check that returns an error if the length of the array is not 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (av *ArrayValidator[T]) IsEmpty() *ArrayValidator[T] {
	av.checks.AddIsEmpty()
	return av
}

// IsNotEmpty adds a check that returns an error if the length of the array is
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (av *ArrayValidator[T]) IsNotEmpty() *ArrayValidator[T] {
	av.checks.AddIsNotEmpty()
	return av
}

// HasCount adds a check that returns an error if the length of the array does not equal the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (av *ArrayValidator[T]) HasCount(l int) *ArrayValidator[T] {
	av.checks.AddHasLength(l)
	return av
}

// HasMoreThan adds a check that returns an error if the length of the array is less than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsGreaterThan(l))
func (av *ArrayValidator[T]) HasMoreThan(l int) *ArrayValidator[T] {
	av.checks.AddIsLongerThan(l)
	return av
}

// HasFewerThan adds a check that returns an error if the length of the array is more than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsLessThan(l))
func (av *ArrayValidator[T]) HasFewerThan(l int) *ArrayValidator[T] {
	av.checks.AddIsShorterThan(l)
	return av
}

// Each assigns a Validator to be used for validating array values
func (av *ArrayValidator[T]) Each(ev with.Validator[T]) *ArrayValidator[T] {
	av.checks.AddIterValValidator(ev)
	return av
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a match for the expected type
func (av *ArrayValidator[T]) ValidateUntyped(value any, options ...*with.ValidationOptions) error {
	if err := testType(value, av.typeStr); err != nil {
		return err
	}
	return av.Validate(value.([]T), options...)
}

// Validate applies all checks against an array and returns an error if any fail
func (av *ArrayValidator[T]) Validate(arr []T, options ...*with.ValidationOptions) error {
	return av.checks.Evaluate(arr, getValidationOptions(options))
}

// Is adds the provided function as a check against any values to be validated
func (av *ArrayValidator[T]) Is(fn func([]T) error) *ArrayValidator[T] {
	av.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		return fn(val)
	})
	return av
}

// Has adds the provided function as a check against any values to be validated
// Has is an alias for Is
func (av *ArrayValidator[T]) Has(fn func([]T) error) *ArrayValidator[T] {
	return av.Is(fn)
}

// Contains causes a validation error if the provided value is not in the array
func (cv *ComparableArrayValidator[T]) Contains(item T) *ComparableArrayValidator[T] {
	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		for _, v := range val {
			if v == item {
				return nil
			}
		}

		return fmt.Errorf(`array must contain value "%v"`, item)
	})
	return cv
}

// DoesNotContain causes a validation error if the provided value is not in the array
func (cv *ComparableArrayValidator[T]) DoesNotContain(item T) *ComparableArrayValidator[T] {
	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		for _, v := range val {
			if v == item {
				return fmt.Errorf(`array must not contain value "%v"`, item)
			}
		}

		return nil
	})
	return cv
}

// ContainsOnly causes a validation error if the array contains any value not in the provided list
func (cv *ComparableArrayValidator[T]) ContainsOnly(items []T) *ComparableArrayValidator[T] {
	allow := map[T]bool{}

	for _, item := range items {
		allow[item] = true
	}

	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		for _, v := range val {
			_, ok := allow[v]

			if !ok {
				return fmt.Errorf(`array must contain only allowed values; value "%v" not allowed`, v)
			}
		}

		return nil
	})
	return cv
}

// ContainsNoDuplicates causes a validation error if the array contains any duplicates
func (cv *ComparableArrayValidator[T]) ContainsNoDuplicates() *ComparableArrayValidator[T] {
	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		found := map[T]bool{}

		for _, v := range val {
			_, ok := found[v]

			if ok {
				return fmt.Errorf(`array must not contain duplicate values; value "%v" is repeated`, v)
			}

			found[v] = true
		}

		return nil
	})
	return cv
}

// ContainsAnyOf causes a validation error if at least one of the provided values is not in the array
func (cv *ComparableArrayValidator[T]) ContainsAnyOf(items []T) *ComparableArrayValidator[T] {
	expect := map[T]bool{}

	for _, item := range items {
		expect[item] = true
	}

	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		for _, v := range val {
			_, ok := expect[v]

			if ok {
				return nil
			}
		}

		return fmt.Errorf(`array must contain one of the expected values`)
	})
	return cv
}

// DoesNotContainAnyOf causes a validation error if at least one of the provided values is in the array
func (cv *ComparableArrayValidator[T]) DoesNotContainAnyOf(items []T) *ComparableArrayValidator[T] {
	expect := map[T]bool{}

	for _, item := range items {
		expect[item] = true
	}

	cv.checks.Append(func(val []T, _ *with.ValidationOptions) error {
		for _, v := range val {
			_, ok := expect[v]

			if ok {
				return fmt.Errorf(`array must not contain any prohibited values`)
			}
		}

		return nil
	})
	return cv
}
