package ensure

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"reflect"
	"strings"
)

type NumberType interface {
	constraints.Integer | constraints.Float
}

type numCheckFunc[T NumberType] func(T) error

// isEven returns a boolean value indicating whether the provided number is even
func isEven[T NumberType](typeStr string, i T) bool {
	// always either 1 or 0, so we don't need a wide int
	// setting result to 1 means the default case is false
	result := int8(1)

	switch typeStr {
	case "int":
		result = int8(int(i) & 1)
	case "uint":
		result = int8(uint(i) & 1)
	case "int8":
		result = int8(i) & 1
	case "uint8":
		result = int8(uint8(i) & 1)
	case "int16":
		result = int8(int16(i) & 1)
	case "uint16":
		result = int8(uint16(i) & 1)
	case "int32":
		result = int8(int32(i) & 1)
	case "uint32":
		result = int8(uint32(i) & 1)
	case "int64":
		result = int8(int64(i) & 1)
	case "uint64":
		result = int8(uint64(i) & 1)
	case "float32":
		fallthrough
	case "float64":
		// can only be even if there is no decimal
		if math.Mod(float64(i), 1) != 0 {
			return false
		}
		result = int8(int64(i) & 1)
	default:
		panic(fmt.Sprintf(`type "%s" cannot be even or odd`, typeStr))
	}

	// note that the default case is to just return false
	return result == 0
}

type NumberValidator[T NumberType] struct {
	typeStr     string
	isFloat     bool
	checks      []numCheckFunc[T]
	placeholder string
}

func (v *NumberValidator[T]) Type() string {
	return v.typeStr
}

// Number constructs a NumberValidator instance and returns a pointer to it
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

// InRange adds a check that returns an error if number being validated is not between the two numbers provided
// Range is inclusive of the lower bound and exclusive of the upper bound
func (v *NumberValidator[T]) InRange(min T, max T) *NumberValidator[T] {
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
		if !isEven[T](v.typeStr, i) {
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
		if isEven[T](v.typeStr, i) {
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
			return errors.New(
				fmt.Sprintf(`number must be one of the permitted values`),
			)
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
			return errors.New(
				fmt.Sprintf(`number must not be one of the prohibited values`),
			)
		}
		return nil
	})
}

// Validate applies all checks against the value being validated and returns an error if any fail
func (v *NumberValidator[T]) Validate(i interface{}) error {
	if err := testType(i, v.typeStr); err != nil {
		return err
	}

	for _, fn := range v.checks {
		if err := fn(i.(T)); err != nil {
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
