package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
)

const defaultAnyValidatorError = "none of the required validators passed"

type AnyValidator struct {
	validators []with.Validator
	t          string
	err        string
}

func Any(validator with.Validator, additional ...with.Validator) *AnyValidator {
	typeStr := validator.Type()

	validators := []with.Validator{
		validator,
	}

	for _, validator := range additional {
		if typeStr != validator.Type() {
			panic(fmt.Sprintf(
				"all validators must be the same type; expected %s, got %s",
				typeStr,
				validator.Type(),
			))
		}

		validators = append(validators, validator)
	}

	return &AnyValidator{
		validators: validators,
		t:          typeStr,
		err:        defaultAnyValidatorError,
	}
}

func (av *AnyValidator) WithError(str string) *AnyValidator {
	av.err = str
	return av
}

func (av *AnyValidator) Type() string {
	return av.t
}

func (av *AnyValidator) Validate(i any) error {
	for _, validator := range av.validators {
		if err := validator.Validate(i); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return NewValidationError(av.err)
}
