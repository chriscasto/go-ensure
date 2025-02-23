package ensure

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"reflect"
	"strings"
)

type NumberType interface {
	constraints.Integer | constraints.Float
}

type numCheckFunc[T NumberType] func(T) error

type NumberValidator[T NumberType] struct {
	typeStr     string
	checks      []numCheckFunc[T]
	placeholder string
}

func (v *NumberValidator[T]) Type() string {
	return v.typeStr
}

func Number[T constraints.Integer | constraints.Float]() *NumberValidator[T] {
	var zero T

	kind := reflect.TypeOf(zero).Kind()

	// Int placeholder value
	ph := "%d"

	// if it's actually a Float, use that placeholder instead
	if string(kind.String()[0]) == "f" {
		ph = "%g"
	}

	return &NumberValidator[T]{
		//typeStr: reflect.TypeOf(zero).String(),
		typeStr:     reflect.TypeOf(zero).String(),
		placeholder: ph,
	}
}

// Replace occurrences of "{}" in error messages with type-specific placeholder
func (v *NumberValidator[T]) fmtErrorMsg(msg string) string {
	return strings.Replace(msg, "{}", v.placeholder, -1)
}

// InRange returns an error if number being validated is not between the two numbers provided
// Range is inclusive of the lower bound and exclusive of the upper bound
func (v *NumberValidator[T]) InRange(min T, max T) *NumberValidator[T] {
	if max < min {
		panic(fmt.Sprintf("max cannot be less than min"))
	}

	return v.Is(func(i T) error {
		if i < min || i >= max {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is out of range ({}, {})"),
					i, min, max),
			)
		}

		return nil
	})
}

func (v *NumberValidator[T]) IsLessThan(max T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i >= max {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is greater than max ({})"), i, max),
			)
		}

		return nil
	})
}

func (v *NumberValidator[T]) IsGreaterThan(min T) *NumberValidator[T] {
	return v.Is(func(i T) error {
		if i <= min {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is less than min ({})"), i, min),
			)
		}

		return nil
	})
}

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

func (v *NumberValidator[T]) Is(fn numCheckFunc[T]) *NumberValidator[T] {
	v.checks = append(v.checks, fn)
	return v
}
