package ensure

import (
	"errors"
	"fmt"
	"reflect"
)

type MapValidator[K comparable, V any] struct {
	typeStr        string
	tests          []func(map[K]V) error
	keyTypeStr     string
	valueTypeStr   string
	keyValidator   Validator
	valueValidator Validator
}

func Map[K comparable, V any]() *MapValidator[K, V] {
	var mapZero map[K]V
	var keyZero K
	var valZero V

	return &MapValidator[K, V]{
		typeStr:      reflect.TypeOf(mapZero).String(),
		keyTypeStr:   reflect.TypeOf(keyZero).String(),
		valueTypeStr: reflect.TypeOf(valZero).String(),
	}
}

func (mv *MapValidator[K, V]) Type() string {
	return mv.typeStr
}

func (mv *MapValidator[K, V]) Validate(i interface{}) error {
	if err := testType(i, mv.typeStr); err != nil {
		return err
	}

	mp := i.(map[K]V)

	for _, fn := range mv.tests {
		if err := fn(mp); err != nil {
			return NewValidationError(err.Error())
		}
	}

	for key, val := range mp {
		if mv.keyValidator != nil {
			if err := mv.keyValidator.Validate(key); err != nil {
				return err
			}
		}

		if mv.valueValidator != nil {
			if err := mv.valueValidator.Validate(val); err != nil {
				return err
			}
		}
	}

	return nil
}

// EachKey assigns a Validator to be used for validating map keys
func (mv *MapValidator[K, V]) EachKey(kv Validator) *MapValidator[K, V] {
	if mv.keyTypeStr != kv.Type() {
		panic(fmt.Sprintf(
			`map validator has keys with type \"%s\", got key validator with type \"%s\"`,
			mv.keyTypeStr,
			kv.Type(),
		))
	}

	mv.keyValidator = kv
	return mv
}

// EachValue assigns a Validator to be used for validating map values
func (mv *MapValidator[K, V]) EachValue(vv Validator) *MapValidator[K, V] {
	if mv.valueTypeStr != vv.Type() {
		panic(fmt.Sprintf(
			`map validator has values with type \"%s\", got value validator with type \"%s\"`,
			mv.valueTypeStr,
			vv.Type(),
		))
	}

	mv.valueValidator = vv
	return mv
}

func (mv *MapValidator[K, V]) IsNotEmpty() *MapValidator[K, V] {
	mv.tests = append(mv.tests, func(mapVal map[K]V) error {
		if len(mapVal) == 0 {
			return errors.New(
				fmt.Sprintf(`map must not be empty`),
			)
		}

		return nil
	})

	return mv
}

func (mv *MapValidator[K, V]) HasCount(l int) *MapValidator[K, V] {
	mv.tests = append(mv.tests, func(mapVal map[K]V) error {
		if len(mapVal) != l {
			return errors.New(
				fmt.Sprintf(
					`map length must equal %d; got %d`,
					l,
					len(mapVal)),
			)
		}

		return nil
	})

	return mv
}

func (mv *MapValidator[K, V]) HasMoreThan(l int) *MapValidator[K, V] {
	mv.tests = append(mv.tests, func(mapVal map[K]V) error {
		if len(mapVal) <= l {
			return errors.New(
				fmt.Sprintf(
					`map must have a length longer than %d; got %d`,
					l,
					len(mapVal)),
			)
		}

		return nil
	})

	return mv
}

func (mv *MapValidator[K, V]) HasFewerThan(l int) *MapValidator[K, V] {
	mv.tests = append(mv.tests, func(mapVal map[K]V) error {
		if len(mapVal) >= l {
			return errors.New(
				fmt.Sprintf(
					`map must have a length less than %d; got %d`,
					l,
					len(mapVal)),
			)
		}

		return nil
	})

	return mv
}
