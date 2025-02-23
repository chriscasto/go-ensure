package ensure

import (
	"fmt"
	"reflect"
)

type StructValidator[T any] struct {
	zeroVal    T
	refVal     reflect.Value
	validators map[string]Validator
	aliases    FriendlyNames
}

func Struct[T any](fv map[string]Validator, friendlyNames ...FriendlyNames) *StructValidator[T] {
	// Create an empty instance of the struct
	var zero T

	ref := reflect.ValueOf(zero)

	if ref.Kind() != reflect.Struct {
		panic("type T must be a struct")
	}

	// make sure that validator type matches field type
	for name, v := range fv {
		field := ref.FieldByName(name)

		if !field.IsValid() {
			panic(
				fmt.Sprintf(
					"field %s does not exist in struct %s",
					name,
					ref.Type().String(),
				),
			)
		}

		if field.Type().String() != v.Type() {
			panic(
				fmt.Sprintf(
					"field %s is type %s but validator expects %s",
					name,
					field.Type().String(),
					v.Type(),
				),
			)
		}
	}

	aliases := FriendlyNames{}

	// add friendly names for struct fields, if any
	if len(friendlyNames) > 0 {
		for _, names := range friendlyNames {
			for field, alias := range names {
				if _, ok := fv[field]; !ok {
					panic(fmt.Sprintf(`cannot set alias for field "%s"; field does not exist`, field))
				}

				aliases[field] = alias
			}
		}
	}

	return &StructValidator[T]{
		zeroVal:    zero,
		refVal:     ref,
		validators: fv,
		aliases:    aliases,
	}
}

func (v *StructValidator[T]) Type() string {
	return v.refVal.Type().String()
}

func (v *StructValidator[T]) Validate(s interface{}) error {
	sRef := reflect.ValueOf(s)
	sRefType := sRef.Type()

	if !sRef.IsValid() || sRefType != v.refVal.Type() {
		return newTypeErrorFromTypes(v.refVal.Type().String(), sRefType.String())
	}

	return v.ValidateStruct(s.(T))
}

func (v *StructValidator[T]) ValidateStruct(s T) error {
	sRef := reflect.ValueOf(s)

	for field, val := range v.validators {
		fieldVal := sRef.FieldByName(field)
		if err := val.Validate(fieldVal.Interface()); err != nil {
			var printName string

			if alias, ok := v.aliases[field]; ok {
				printName = alias
			} else {
				printName = field
			}

			return NewValidationError(fmt.Sprintf("%s: %s", printName, err.Error()))
		}
	}
	return nil
}
