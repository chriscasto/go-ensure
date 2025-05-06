package ensure

import (
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

const defaultAnyValidatorError = "none of the required validators passed"

type AnyValidator[T any] struct {
	validators []with.Validator[T]
	t          string
	//err        string
	opts *with.AnyOptions
}

// Any instantiates and returns an instance of AnyValidator
func Any[T any](validators ...with.Validator[T]) *AnyValidator[T] {
	var zero T

	typeStr := reflect.ValueOf(zero).Type().String()

	return &AnyValidator[T]{
		validators: validators,
		t:          typeStr,
		//err:        defaultAnyValidatorError,
		opts: with.DefaultAnyOptions(),
	}
}

// WithError sets the default error string to return if none of the validators pass
func (av *AnyValidator[T]) WithError(str string) *AnyValidator[T] {
	with.AnyOptionDefaultError(str)(av.opts)
	return av
}

// Type returns a string with the type this validator expects
func (av *AnyValidator[T]) Type() string {
	return av.t
}

func (av *AnyValidator[T]) WithOptions(opts ...with.AnyOption) *AnyValidator[T] {
	for _, opt := range opts {
		opt(av.opts)
	}
	return av
}

// Validate applies all validators against a value of the expected type and returns an error if all fail
func (av *AnyValidator[T]) Validate(i T, options ...*with.ValidationOptions) error {
	vOpts := getValidationOptions(options)
	errByIdx := make(map[int]error)

	for idx, validator := range av.validators {
		if err := validator.Validate(i, vOpts); err != nil {
			errByIdx[idx] = err
		} else {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, we should have at least one error
	if len(errByIdx) > 0 {
		if vOpts.CollectAllErrors() {
			vErrs := newValidationErrors()

			for idx, err := range errByIdx {
				if av.opts.PassThroughErrorsFrom(idx) {
					vErrs.Append(err)
				}
			}

			if !vErrs.HasErrors() {
				vErrs.Append(av.opts.DefaultError())
			}

			return vErrs
		} else {
			for idx, err := range errByIdx {
				if av.opts.PassThroughErrorsFrom(idx) {
					return err
				}
			}
		}
	}

	// Since we return on the first nil above, the only way we get here is
	// if we aren't passing through any of the errors. We return the default
	// error message instead
	return av.opts.DefaultError()
}

// ValidateUntyped applies all validators against a value of an unknown type and returns an error if all fail
func (av *AnyValidator[T]) ValidateUntyped(i any, options ...*with.ValidationOptions) error {
	for _, validator := range av.validators {
		if err := validator.ValidateUntyped(i, options...); err == nil {
			// If any pass without error, consider it a success
			return nil
		}
	}

	// If we haven't encountered a success, return an error
	return av.opts.DefaultError()
}
