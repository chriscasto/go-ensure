package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

type PointerValidator[T any] struct {
	parent   with.Validator[T]
	optional bool
	t        string
}

func newPtrValidator[T any](parent with.Validator[T], optional bool) *PointerValidator[T] {
	return &PointerValidator[T]{
		parent:   parent,
		optional: optional,
		t:        fmt.Sprintf("*%s", parent.Type()),
	}
}

func Pointer[T any](parent with.Validator[T]) *PointerValidator[T] {
	return newPtrValidator[T](parent, false)
}

func OptionalPointer[T any](parent with.Validator[T]) *PointerValidator[T] {
	return newPtrValidator[T](parent, true)
}

func (v *PointerValidator[T]) Type() string {
	return v.t
}

func (v *PointerValidator[T]) ValidateUntyped(i any) error {
	refVal := reflect.ValueOf(i)

	if refVal.Kind() != reflect.Ptr {
		return NewTypeError("value must be a pointer")
	}

	if refVal.IsNil() {
		if !v.optional {
			return NewValidationError("required value cannot be missing")
		}
		return nil
	}

	return v.parent.Validate(refVal.Elem().Interface().(T))
}

func (v *PointerValidator[T]) Validate(i *T) error {
	if i == nil {
		if !v.optional {
			return NewValidationError("required value cannot be missing")
		}
		return nil
	}

	return v.parent.Validate(*i)
}
