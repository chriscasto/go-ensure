package with

type ValidationOptions struct {
	collectAllErrors bool
}

func (vo *ValidationOptions) CollectAllErrors() bool {
	return vo.collectAllErrors
}

type Option func(*ValidationOptions)

func OptionCollectAllErrors() Option {
	return func(o *ValidationOptions) {
		o.collectAllErrors = true
	}
}

func Options(options ...Option) *ValidationOptions {
	valOpts := &ValidationOptions{
		// defaults
		collectAllErrors: false,
	}

	for _, opt := range options {
		opt(valOpts)
	}

	return valOpts
}
