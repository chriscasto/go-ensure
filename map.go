package ensure

import (
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

// MapValidator contains information and logic used to validate a map with keys of type K and values of type V
type MapValidator[K comparable, V any] struct {
	typeStr      string
	keyTypeStr   string
	valueTypeStr string
	checks       *iterChecks[K, V, map[K]V]
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
		checks:       newMapIterChecks[K, V](),
	}
}

// Type returns a string with the type of map this validator expects
func (mv *MapValidator[K, V]) Type() string {
	return mv.typeStr
}

// HasLengthWhere adds a NumberValidator for validating the length of the string
func (mv *MapValidator[K, V]) HasLengthWhere(nv *NumberValidator[int]) *MapValidator[K, V] {
	mv.checks.AddHasLengthWhere(nv)
	return mv
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a match for the expected type
func (mv *MapValidator[K, V]) ValidateUntyped(value any, options ...*with.ValidationOptions) error {
	if err := testType(value, mv.typeStr); err != nil {
		return err
	}
	return mv.Validate(value.(map[K]V), options...)
}

// Validate applies all checks against a map and returns an error if any fail
func (mv *MapValidator[K, V]) Validate(mp map[K]V, options ...*with.ValidationOptions) error {
	return mv.checks.Evaluate(mp, getValidationOptions(options))
}

// EachKey assigns a Validator to be used for validating map keys
func (mv *MapValidator[K, V]) EachKey(kv with.Validator[K]) *MapValidator[K, V] {
	mv.checks.AddIterKeyValidator(kv)
	return mv
}

// EachValue assigns a Validator to be used for validating map values
func (mv *MapValidator[K, V]) EachValue(vv with.Validator[V]) *MapValidator[K, V] {
	mv.checks.AddIterValValidator(vv)
	return mv
}

// IsEmpty adds a check that returns an error if the length of the map is not 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (mv *MapValidator[K, V]) IsEmpty() *MapValidator[K, V] {
	mv.checks.AddIsEmpty()
	return mv
}

// IsNotEmpty adds a check that returns an error if the length of the map is 0
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (mv *MapValidator[K, V]) IsNotEmpty() *MapValidator[K, V] {
	mv.checks.AddIsNotEmpty()
	return mv
}

// HasCount adds a check that returns an error if the length of the map is not the passed value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (mv *MapValidator[K, V]) HasCount(l int) *MapValidator[K, V] {
	mv.checks.AddHasLength(l)
	return mv
}

// HasMoreThan adds a check that returns an error if the length of the map is less than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsGreaterThan(l))
func (mv *MapValidator[K, V]) HasMoreThan(l int) *MapValidator[K, V] {
	mv.checks.AddIsLongerThan(l)
	return mv
}

// HasFewerThan adds a check that returns an error if the length of the map is more than the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().IsLessThan(l))
func (mv *MapValidator[K, V]) HasFewerThan(l int) *MapValidator[K, V] {
	mv.checks.AddIsShorterThan(l)
	return mv
}

// Is adds the provided function as a check against any values to be validated
func (mv *MapValidator[K, V]) Is(fn func(map[K]V) error) *MapValidator[K, V] {
	mv.checks.Append(func(val map[K]V, _ *with.ValidationOptions) error {
		return fn(val)
	})
	return mv
}
