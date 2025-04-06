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

// newPtrValidator instantiates a new PointerValidator
func newPtrValidator[T any](parent with.Validator[T], optional bool) *PointerValidator[T] {
	return &PointerValidator[T]{
		parent:   parent,
		optional: optional,
		t:        fmt.Sprintf("*%s", parent.Type()),
	}
}

// Pointer returns a PointerValidator that returns an error on a nil pointer
func Pointer[T any](parent with.Validator[T]) *PointerValidator[T] {
	return newPtrValidator[T](parent, false)
}

// OptionalPointer returns a PointerValidator that doesn't return an error on a nil pointer
func OptionalPointer[T any](parent with.Validator[T]) *PointerValidator[T] {
	return newPtrValidator[T](parent, true)
}

// Type returns a string with indicating a pointer to the parent validator type
func (v *PointerValidator[T]) Type() string {
	return v.t
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a pointer to the expected type
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

// Validate applies all checks against a boolean value and returns an error if any fail
func (v *PointerValidator[T]) Validate(i *T) error {
	if i == nil {
		if !v.optional {
			return NewValidationError("required value cannot be missing")
		}
		return nil
	}

	return v.parent.Validate(*i)
}
