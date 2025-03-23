package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

// mapCheckFunc defines a function that can be used to validate a map
type mapCheckFunc[K comparable, V any] func(map[K]V) error

// MapValidator contains information and logic used to validate a map with keys of type K and values of type V
type MapValidator[K comparable, V any] struct {
	typeStr        string
	keyTypeStr     string
	valueTypeStr   string
	checks         []mapCheckFunc[K, V]
	lenValidator   *NumberValidator[int]
	keyValidator   with.Validator
	valueValidator with.Validator
}

// Map constructs a MapValidator instance with keys of type K and values of type V and returns a pointer to it
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

// Type returns a string with the type of map this validator expects
func (mv *MapValidator[K, V]) Type() string {
	return mv.typeStr
}

// HasLengthWhere adds a NumberValidator for validating the length of the string
func (mv *MapValidator[K, V]) HasLengthWhere(nv *NumberValidator[int]) *MapValidator[K, V] {
	mv.lenValidator = nv
	return mv
}

// Validate accepts an arbitrary input type and validates it if it's a match for the expected type
func (mv *MapValidator[K, V]) Validate(i interface{}) error {
	if err := testType(i, mv.typeStr); err != nil {
		return err
	}
	return mv.ValidateMap(i.(map[K]V))
}

// ValidateMap applies all checks against a map and returns an error if any fail
func (mv *MapValidator[K, V]) ValidateMap(mp map[K]V) error {
	if mv.lenValidator != nil {
		if err := mv.lenValidator.Validate(len(mp)); err != nil {
			return err
		}
	}

	for _, fn := range mv.checks {
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
func (mv *MapValidator[K, V]) EachKey(kv with.Validator) *MapValidator[K, V] {
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
func (mv *MapValidator[K, V]) EachValue(vv with.Validator) *MapValidator[K, V] {
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

// IsEmpty adds a check that returns an error if the length of the map is not 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (mv *MapValidator[K, V]) IsEmpty() *MapValidator[K, V] {
	return mv.Is(func(mapVal map[K]V) error {
		if len(mapVal) != 0 {
			return errors.New(`map must be empty`)
		}

		return nil
	})
}

// IsNotEmpty adds a check that returns an error if the length of the map is 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (mv *MapValidator[K, V]) IsNotEmpty() *MapValidator[K, V] {
	return mv.Is(func(mapVal map[K]V) error {
		if len(mapVal) == 0 {
			return errors.New(`map must not be empty`)
		}

		return nil
	})
}

// HasCount adds a check that returns an error if the length of the map is not the passed value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (mv *MapValidator[K, V]) HasCount(l int) *MapValidator[K, V] {
	return mv.Is(func(mapVal map[K]V) error {
		if len(mapVal) != l {
			return errors.New(
				fmt.Sprintf(`map length must equal %d; got %d`, l, len(mapVal)),
			)
		}

		return nil
	})
}

// HasMoreThan adds a check that returns an error if the length of the map is less than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsGreaterThan(l))
func (mv *MapValidator[K, V]) HasMoreThan(l int) *MapValidator[K, V] {
	return mv.Is(func(mapVal map[K]V) error {
		if len(mapVal) <= l {
			return errors.New(
				fmt.Sprintf(`map must have a length longer than %d; got %d`, l, len(mapVal)),
			)
		}

		return nil
	})
}

// HasFewerThan adds a check that returns an error if the length of the map is more than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsLessThan(l))
func (mv *MapValidator[K, V]) HasFewerThan(l int) *MapValidator[K, V] {
	return mv.Is(func(mapVal map[K]V) error {
		if len(mapVal) >= l {
			return errors.New(
				fmt.Sprintf(`map must have a length less than %d; got %d`, l, len(mapVal)),
			)
		}

		return nil
	})
}

// Is adds the provided function as a check against any values to be validated
func (mv *MapValidator[K, V]) Is(fn mapCheckFunc[K, V]) *MapValidator[K, V] {
	mv.checks = append(mv.checks, fn)
	return mv
}
