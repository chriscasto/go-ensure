package ensure

import (
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

const defaultAnyValidatorError = "none of the required validators passed"

type AnyValidator[T any] struct {
	validators []with.TypedValidator[T]
	t          string
	err        string
}

func Any[T any](validators ...with.TypedValidator[T]) *AnyValidator[T] {
	var zero T

	typeStr := reflect.ValueOf(zero).Type().String()

	return &AnyValidator[T]{
		validators: validators,
		t:          typeStr,
		err:        defaultAnyValidatorError,
	}
}

func (av *AnyValidator[T]) WithError(str string) *AnyValidator[T] {
	av.err = str
	return av
}

func (av *AnyValidator[T]) Type() string {
	return av.t
}

func (av *AnyValidator[T]) ValidateStrict(i T) error {
	for _, validator := range av.validators {
		if err := validator.ValidateStrict(i); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return NewValidationError(av.err)
}

func (av *AnyValidator[T]) Validate(i any) error {
	for _, validator := range av.validators {
		if err := validator.Validate(i); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return NewValidationError(av.err)
}
