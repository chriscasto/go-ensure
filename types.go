package ensure

type Validator interface {
	Validate(interface{}) error
	Type() string
}

type Fields map[string]Validator
