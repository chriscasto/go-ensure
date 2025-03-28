package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

type PointerValidator struct {
	parent   with.Validator
	optional bool
	t        string
}

func newPtrValidator(parent with.Validator, optional bool) *PointerValidator {
	return &PointerValidator{
		parent:   parent,
		optional: optional,
		t:        fmt.Sprintf("*%s", parent.Type()),
	}
}

func Pointer(parent with.Validator) *PointerValidator {
	return newPtrValidator(parent, false)
}

func OptionalPointer(parent with.Validator) *PointerValidator {
	return newPtrValidator(parent, true)
}

func (v *PointerValidator) Type() string {
	return v.t
}

func (v *PointerValidator) Validate(i any) error {
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

	return v.parent.Validate(refVal.Elem().Interface())
}
