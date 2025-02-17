package valid

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

type NumberValidator[T NumberType] struct {
	//zeroVal T
	typeStr     string
	kind        reflect.Kind
	tests       []func(T) error
	placeholder string
}

func (v *NumberValidator[T]) Type() string {
	return v.typeStr
}

func (v *NumberValidator[T]) Kind() reflect.Kind {
	return v.kind
}

func Number[T constraints.Integer | constraints.Float]() *NumberValidator[T] {
	var zero T

	kind := reflect.TypeOf(zero).Kind()

	// int placeholder value
	ph := "%d"

	// if it's actually a float, use that placeholder instead
	if string(kind.String()[0]) == "f" {
		ph = "%g"
	}

	return &NumberValidator[T]{
		//typeStr: reflect.TypeOf(zero).String(),
		typeStr:     reflect.TypeOf(zero).String(),
		kind:        kind,
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

	v.tests = append(v.tests, func(i T) error {
		if i < min || i >= max {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is out of range ({}, {})"),
					i, min, max),
			)
		}

		return nil
	})

	return v
}

func (v *NumberValidator[T]) LessThan(max T) *NumberValidator[T] {
	v.tests = append(v.tests, func(i T) error {
		if i >= max {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is greater than max ({})"), i, max),
			)
		}

		return nil
	})

	return v
}

func (v *NumberValidator[T]) GreaterThan(min T) *NumberValidator[T] {
	v.tests = append(v.tests, func(i T) error {
		if i <= min {
			return errors.New(
				fmt.Sprintf(
					v.fmtErrorMsg("number {} is less than min ({})"), i, min),
			)
		}

		return nil
	})

	return v
}

func (v *NumberValidator[T]) Validate(i interface{}) error {
	valKind := reflect.TypeOf(i).Kind()
	if valKind != v.kind {
		return fmt.Errorf(
			`number validator expects type "%s"; got "%s"`,
			v.kind.String(),
			valKind.String(),
		)
	}

	for _, fn := range v.tests {
		err := fn(i.(T))
		if err != nil {
			return err
		}
	}

	return nil
}
