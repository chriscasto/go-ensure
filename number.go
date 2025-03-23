package ensure

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"reflect"
	"strings"
)

// NumberType defines the set of values accepted by NumberValidator
type NumberType interface {
	constraints.Integer | constraints.Float
}

// numCheckFunc is the type of function used for validating number types
type numCheckFunc[T NumberType] func(T) error

// isIntEven returns true if an int value can be considered even
func isIntEven(typeStr string, i any) bool {
	// we only use the last bit, so we don't need a wide int
	var intVal int8

	// coerce to the appropriate type, then downcast to int8
	switch typeStr {
	case "int":
		intVal = int8(i.(int))
	case "uint":
		intVal = int8(i.(uint))
	case "int8":
		intVal, _ = i.(int8)
	case "uint8":
		intVal = int8(i.(uint8))
	case "int16":
		intVal = int8(i.(int16))
	case "uint16":
		intVal = int8(i.(uint16))
	case "int32":
		intVal = int8(i.(int32))
	case "uint32":
		intVal = int8(i.(uint32))
	case "int64":
		intVal = int8(i.(int64))
	case "uint64":
		intVal = int8(i.(uint64))
	default:
		panic(fmt.Sprintf(`type "%s" cannot be even or odd`, typeStr))
	}

	// int is even if the 1 bit is not set
	return (intVal & 1) == 0
}

// isFloatEven returns true if a float value can be considered even
func isFloatEven(i float64) bool {
	// can't be even if the decimal component is not zero
	if math.Mod(i, 1) != 0 {
		return false
	}

	// int component is even if the 1 bit is not set
	return int8(int64(i)&1) == 0
}

// isFloatOdd returns true if a float value can be considered odd
func isFloatOdd(i float64) bool {
	// can't be odd if the decimal component is not zero
	if math.Mod(i, 1) != 0 {
		return false
	}

	// int component is odd if the 1 bit is set
	return int8(int64(i)&1) == 1
}

// isEven returns a boolean value indicating whether the provided number is even
func isEven(typeStr string, i any) bool {
	// check to see whether it's a float
	switch typeStr {
	case "float32":
		return isFloatEven(float64(i.(float32)))
	case "float64":
		return isFloatEven(i.(float64))
	default:
		return isIntEven(typeStr, i)
	}
}

// isOdd returns a boolean value indicating whether the provided number is odd
func isOdd(typeStr string, i any) bool {
	// check to see whether it's a float
	switch typeStr {
	case "float32":
		return isFloatOdd(float64(i.(float32)))
	case "float64":
		return isFloatOdd(i.(float64))
	default:
		return !isIntEven(typeStr, i)
	}
}

// NumberValidator contains information and logic used to validate a number of type T
type NumberValidator[T NumberType] struct {
	typeStr     string
	isFloat     bool
	checks      []numCheckFunc[T]
	placeholder string
}

// Type returns a string with the type of the number this validator expects
func (v *NumberValidator[T]) Type() string {
	return v.typeStr
}

// Number constructs a NumberValidator instance of type T and returns a pointer to it
func Number[T constraints.Integer | constraints.Float]() *NumberValidator[T] {
	var zero T

	kind := reflect.TypeOf(zero).Kind()

	// Int placeholder value
	ph := "%d"
	isFloat := false

	// if it's actually a Float, use that placeholder instead
	if string(kind.String()[0]) == "f" {
		ph = "%g"
		isFloat = true
	}

	return &NumberValidator[T]{
		typeStr:     reflect.TypeOf(zero).String(),
		placeholder: ph,
		isFloat:     isFloat,
	}
}

// fmtErrorMsg replaces occurrences of "{}" in error messages with type-specific placeholder
func (v *NumberValidator[T]) fmtErrorMsg(msg string) string {
	return strings.Replace(msg, "{}", v.placeholder, -1)
}

// IsInRange adds a check that returns an error if number being validated is not between the two numbers provided
// Range is inclusive of the lower bound and exclusive of the upper bound
func (v *NumberValidator[T]) IsInRange(min T, max T) *NumberValidator[T] {
	if max < min {
		panic(fmt.Sprintf("max cannot be less than min"))
	}

	return v.Is(func(i T) error {
		if i < min || i >= max {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be in the range [{}, {}); got {}"),
					min, max, i),
			)
		}

		return nil
	})
}

// Equals adds a check that returns an error if number being validated is not exactly the number provided
func (v *NumberValidator[T]) Equals(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i != target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must equal {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// DoesNotEqual adds a check that returns an error if number being validated is exactly the same as the number provided
func (v *NumberValidator[T]) DoesNotEqual(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i == target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must not equal {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// IsLessThan adds a check that returns an error if number being validated is not lees than the number provided
func (v *NumberValidator[T]) IsLessThan(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i >= target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be less than {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// IsLessThanOrEqualTo adds a check that returns an error if number being validated is not lees than or equal to the number provided
func (v *NumberValidator[T]) IsLessThanOrEqualTo(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i > target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be less than or equal to {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// IsGreaterThan adds a check that returns an error if number being validated is not greater than the number provided
func (v *NumberValidator[T]) IsGreaterThan(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i <= target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be greater than {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// IsGreaterThanOrEqualTo adds a check that returns an error if number being validated is not greater than or equal to than the number provided
func (v *NumberValidator[T]) IsGreaterThanOrEqualTo(target T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i < target {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be greater than or equal to {}; got {}"), target, i),
			)
		}

		return nil
	})
}

// IsEven adds a check that returns an error if number being validated is not even
func (v *NumberValidator[T]) IsEven() *NumberValidator[T] {
	return v.Is(func(i T) error {
		if !isEven(v.typeStr, i) {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be even; got {}"), i),
			)
		}

		return nil
	})
}

// IsOdd adds a check that returns an error if number being validated is not odd
func (v *NumberValidator[T]) IsOdd() *NumberValidator[T] {
	return v.Is(func(i T) error {
		if !isOdd(v.typeStr, i) {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number must be odd; got {}"), i),
			)
		}

		return nil
	})
}

// IsPositive is a shortcut for IsGreaterThan(0)
func (v *NumberValidator[T]) IsPositive() *NumberValidator[T] {
	return v.IsGreaterThan(0)
}

// IsNegative is a shortcut for IsLessThan(0)
func (v *NumberValidator[T]) IsNegative() *NumberValidator[T] {
	return v.IsLessThan(0)
}

// IsZero is a shortcut for Equals(0)
func (v *NumberValidator[T]) IsZero() *NumberValidator[T] {
	return v.Equals(0)
}

// IsNotZero is a shortcut for DoesNotEqual(0)
func (v *NumberValidator[T]) IsNotZero() *NumberValidator[T] {
	return v.DoesNotEqual(0)
}

// IsOneOf adds a check that returns an error if number being validated is not in the provided list
func (v *NumberValidator[T]) IsOneOf(values []T) *NumberValidator[T] {
	// convert list to map for O(1) lookups
	lookup := map[T]bool{}

	for _, num := range values {
		lookup[num] = true
	}

	return v.Is(func(num T) error {
		if _, ok := lookup[num]; !ok {
			return errors.New(`number must be one of the permitted values`)
		}
		return nil
	})
}

// IsNotOneOf adds a check that returns an error if number being validated is in the provided list
func (v *NumberValidator[T]) IsNotOneOf(values []T) *NumberValidator[T] {
	// convert list to map for O(1) lookups
	lookup := map[T]bool{}

	for _, num := range values {
		lookup[num] = true
	}

	return v.Is(func(num T) error {
		if _, ok := lookup[num]; ok {
			return errors.New(`number must not be one of the prohibited values`)
		}
		return nil
	})
}

// Validate accepts an arbitrary input type and validates it if it's a match for the expected type
func (v *NumberValidator[T]) Validate(i interface{}) error {
	if err := testType(i, v.typeStr); err != nil {
		return err
	}

	return v.ValidateNumber(i.(T))
}

// ValidateNumber applies all checks against a number of the expected type and returns an error if any fail
func (v *NumberValidator[T]) ValidateNumber(n T) error {
	for _, fn := range v.checks {
		if err := fn(n); err != nil {
			return err
		}
	}

	return nil
}

// Is adds the provided function as a check against any values to be validated
func (v *NumberValidator[T]) Is(fn numCheckFunc[T]) *NumberValidator[T] {
	v.checks = append(v.checks, fn)
	return v
}
