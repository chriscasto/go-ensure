package ensure

import (
	"errors"
	"fmt"
	"reflect"
)

type arrCheckFunc[T any] func([]T) error

type ArrayValidator[T any] struct {
	typeStr string
	checks  []arrCheckFunc[T]
}

func Array[T any]() *ArrayValidator[T] {
	var zero T

	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())

	return &ArrayValidator[T]{
		typeStr: typeStr,
	}
}

func (v *ArrayValidator[T]) Type() string {
	return v.typeStr
}

func (v *ArrayValidator[T]) IsNotEmpty() *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		if len(arr) == 0 {
			return errors.New(
				fmt.Sprintf(`array must not be empty`),
			)
		}

		return nil
	})
}

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

func (v *ArrayValidator[T]) Each(ev Validator) *ArrayValidator[T] {
	return v.Is(func(arr []T) error {
		for _, e := range arr {
			if err := ev.Validate(e); err != nil {
				return err
			}
		}

		return nil
	})
}

func (v *ArrayValidator[T]) Validate(i interface{}) error {
	if err := testType(i, v.typeStr); err != nil {
		return err
	}

	for _, fn := range v.checks {
		if err := fn(i.([]T)); err != nil {
			return err
		}
	}

	return nil
}

func (v *ArrayValidator[T]) Is(fn arrCheckFunc[T]) *ArrayValidator[T] {
	v.checks = append(v.checks, fn)
	return v
}
