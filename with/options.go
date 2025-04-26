package with

// ValidationOptions is a struct containing all settings for performing validation
type ValidationOptions struct {
	collectAllErrors bool
}

// CollectAllErrors returns true if all checks need to be evaluated and all errors returned collected
// A value of false means that validation should return immediately on the first error
func (vo *ValidationOptions) CollectAllErrors() bool {
	return vo.collectAllErrors
}

// ValidationOption is a function signature for an option that can be applied to validation settings
type ValidationOption func(*ValidationOptions)

// OptionCollectAllErrors causes CollectAllErrors() to return true
func OptionCollectAllErrors() ValidationOption {
	return func(o *ValidationOptions) {
		o.collectAllErrors = true
	}
}

// DefaultValidationOptions returns ValidationOptions with the default values set
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		collectAllErrors: false,
	}
}

// Options returns a ValidationOptions struct with all passed options applied
func Options(options ...ValidationOption) *ValidationOptions {
	valOpts := DefaultValidationOptions()

	for _, opt := range options {
		opt(valOpts)
	}

	return valOpts
}
