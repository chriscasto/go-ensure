package valid

import (
	"errors"
	"fmt"
	"reflect"
)

type ArrayValidator[T any] struct {
	kind    reflect.Kind
	typeStr string
	tests   []func([]T) error
}

//func (v *ArrayValidator) Valid(arr []interface{}) bool {}

func Array[T any]() *ArrayValidator[T] {
	var zero T

	//kind := reflect.TypeOf(zero).Kind()

	//fmt.Println(reflect.TypeOf(zero).String())
	typeStr := fmt.Sprintf("[]%s", reflect.TypeOf(zero).String())
	//fmt.Println(typeStr)

	//fmt.Println(kind.String())
	return &ArrayValidator[T]{
		kind:    reflect.Slice,
		typeStr: typeStr,
	}
}

func (v *ArrayValidator[T]) Type() string {
	return v.typeStr
}

func (v *ArrayValidator[T]) Kind() reflect.Kind {
	return v.kind
}

func (v *ArrayValidator[T]) NotEmpty() *ArrayValidator[T] {
	v.tests = append(v.tests, func(arr []T) error {
		if len(arr) == 0 {
			return errors.New(
				fmt.Sprintf(`array is empty but shouldn't be`),
			)
		}

		return nil
	})

	return v
}

func (v *ArrayValidator[T]) Count(l int) *ArrayValidator[T] {
	v.tests = append(v.tests, func(arr []T) error {
		if len(arr) != l {
			return errors.New(
				fmt.Sprintf(
					`array has length %d but expects %d`,
					len(arr),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *ArrayValidator[T]) MoreThan(l int) *ArrayValidator[T] {
	v.tests = append(v.tests, func(arr []T) error {
		if len(arr) <= l {
			return errors.New(
				fmt.Sprintf(
					`array has length %d which is shorter than min (%d)`,
					len(arr),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *ArrayValidator[T]) FewerThan(l int) *ArrayValidator[T] {
	v.tests = append(v.tests, func(arr []T) error {
		if len(arr) >= l {
			return errors.New(
				fmt.Sprintf(
					`array has length %d which is longer than max (%d)`,
					len(arr),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *ArrayValidator[T]) Each(ev Validator) *ArrayValidator[T] {
	v.tests = append(v.tests, func(arr []T) error {
		for _, e := range arr {
			if err := ev.Validate(e); err != nil {
				return err
			}
		}

		return nil
	})

	return v
}

func (v *ArrayValidator[T]) Validate(i interface{}) error {
	valType := reflect.TypeOf(i).String()

	if valType != v.typeStr {
		return fmt.Errorf(
			`array validator expects type "%s"; got "%s"`,
			v.typeStr,
			valType,
		)
	}

	for _, fn := range v.tests {
		err := fn(i.([]T))
		if err != nil {
			return err
		}
	}

	return nil
}
