package with

import "errors"

const defaultAnyValidatorError = "none of the required validators passed"

// AnyOptions are a set of options that apply to a specific Any validator
type AnyOptions struct {
	defaultErr string
	passThru   map[int]bool
}

func DefaultAnyOptions() *AnyOptions {
	return &AnyOptions{
		defaultErr: defaultAnyValidatorError,
		passThru:   make(map[int]bool),
	}
}

func (opts *AnyOptions) DefaultError() error {
	return errors.New(opts.defaultErr)
}

// PassThroughErrorsFrom indicates whether the validator should pass through
// errors received from a validators at a specific index
func (opts *AnyOptions) PassThroughErrorsFrom(idx int) bool {
	pass, ok := opts.passThru[idx]

	if ok && pass {
		return true
	}

	return false
}

type AnyOption func(*AnyOptions)

func AnyOptionDefaultError(msg string) AnyOption {
	return func(o *AnyOptions) {
		o.defaultErr = msg
	}
}

func AnyOptionPassThroughErrorsFrom(idx ...int) AnyOption {
	return func(o *AnyOptions) {
		for _, i := range idx {
			o.passThru[i] = true
		}
	}
}
