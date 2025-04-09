package with

type ValidationOptions struct {
}

type Option func(*ValidationOptions)

func Options(options ...Option) *ValidationOptions {
	valOpts := &ValidationOptions{}

	for _, opt := range options {
		opt(valOpts)
	}

	return valOpts
}
