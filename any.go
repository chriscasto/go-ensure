package ensure

import (
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

const defaultAnyValidatorError = "none of the required validators passed"

type AnyValidator[T any] struct {
	validators []with.Validator[T]
	t          string
	err        string
}

// Any instantiates and returns an instance of AnyValidator
func Any[T any](validators ...with.Validator[T]) *AnyValidator[T] {
	var zero T

	typeStr := reflect.ValueOf(zero).Type().String()

	return &AnyValidator[T]{
		validators: validators,
		t:          typeStr,
		err:        defaultAnyValidatorError,
	}
}

// WithError sets the default error string to return if none of the validators pass
func (av *AnyValidator[T]) WithError(str string) *AnyValidator[T] {
	av.err = str
	return av
}

// Type returns a string with the type this validator expects
func (av *AnyValidator[T]) Type() string {
	return av.t
}

// Validate applies all validators against a value of the expected type and returns an error if all fail
func (av *AnyValidator[T]) Validate(i T) error {
	for _, validator := range av.validators {
		if err := validator.Validate(i); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return NewValidationError(av.err)
}

// ValidateUntyped applies all validators against a value of an unknown type and returns an error if all fail
func (av *AnyValidator[T]) ValidateUntyped(i any) error {
	for _, validator := range av.validators {
		if err := validator.ValidateUntyped(i); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return NewValidationError(av.err)
}
