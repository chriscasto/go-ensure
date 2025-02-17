package valid

import "reflect"

type Validator interface {
	Validate(interface{}) error
	Type() string
	Kind() reflect.Kind
}

type Fields map[string]Validator
