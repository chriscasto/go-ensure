package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

const defaultAnyValidatorError = "none of the required validators passed"

type AnyValidator[T any] struct {
	validators []with.Validator
	t          string
	err        string
}

func Any[T any](validators ...with.Validator) *AnyValidator[T] {
	var zero T

	typeStr := reflect.ValueOf(zero).Type().String()
	//validator.Type()

	//validators := []with.Validator{
	//	validator,
	//}

	// Check to make sure that each validator implements the right type
	for _, validator := range validators {
		if typeStr != validator.Type() {
			panic(fmt.Sprintf(
				"all validators must be the same type; expected %s, got %s",
				typeStr,
				validator.Type(),
			))
		}

		//validators = append(validators, validator)
	}

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

func (av *AnyValidator[T]) ValidateStrict(t T) error {
	for _, validator := range av.validators {
		strict, ok := validator.(with.StrictValidator[T])
		var err error

		// First try to use strict validation
		if ok {
			err = strict.ValidateStrict(t)
		} else {
			err = validator.Validate(t)
		}

		if err == nil {
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
