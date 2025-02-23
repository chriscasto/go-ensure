package ensure

type Validator interface {
	Validate(interface{}) error
	Type() string
}

type Fields map[string]Validator

type TypeError struct {
	err string
}

func NewTypeError(err string) *TypeError {
	return &TypeError{err}
}

func (e *TypeError) Error() string {
	return e.err
}
