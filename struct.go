package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

type structCheckFunc[T any] func(T) error

type StructValidator[T any] struct {
	zeroVal T
	refVal  reflect.Value
	tests   []structCheckFunc[T]

	fieldValidators with.Validators
	fieldAliases    with.FriendlyNames

	getterValidators with.Validators
	getterAliases    with.FriendlyNames
}

func Struct[T any]() *StructValidator[T] {
	// Create an empty instance of the struct
	var zero T

	ref := reflect.ValueOf(zero)

	if ref.Kind() != reflect.Struct {
		panic("type T must be a struct")
	}

	return &StructValidator[T]{
		zeroVal:          zero,
		refVal:           ref,
		fieldValidators:  with.Validators{},
		fieldAliases:     with.FriendlyNames{},
		getterValidators: with.Validators{},
		getterAliases:    with.FriendlyNames{},
	}
}

func (sv *StructValidator[T]) HasFields(validators with.Validators, friendlyNames ...with.FriendlyNames) *StructValidator[T] {
	ref := sv.refVal

	// make sure that validator type matches field type
	for name, validator := range validators {
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

		if field.Type().String() != validator.Type() {
			panic(
				fmt.Sprintf(
					"field %s is type [%s] but validator expects [%s]",
					name,
					field.Type().String(),
					validator.Type(),
				),
			)
		}

		sv.fieldValidators[name] = validator
	}

	// add friendly names for struct validators, if any
	if len(friendlyNames) > 0 {
		for _, names := range friendlyNames {
			for field, alias := range names {
				if _, ok := validators[field]; !ok {
					panic(fmt.Sprintf(`cannot set alias for field "%s"; field does not exist`, field))
				}

				sv.fieldAliases[field] = alias
			}
		}
	}

	return sv
}

func (sv *StructValidator[T]) HasGetters(validators with.Validators, friendlyNames ...with.FriendlyNames) *StructValidator[T] {
	ref := sv.refVal

	// make sure that validator type matches field type
	for name, validator := range validators {
		method := ref.MethodByName(name)

		if !method.IsValid() {
			panic(
				fmt.Sprintf(
					"method %s() does not exist in struct %s",
					name,
					ref.Type().String(),
				),
			)
		}

		mType := method.Type()

		if mType.NumOut() > 1 {
			panic(
				fmt.Sprintf(
					"method %s() has %d return values but validator expects 1",
					name,
					mType.NumOut(),
				),
			)
		}

		retVal := mType.Out(0)

		if retVal.String() != validator.Type() {
			panic(
				fmt.Sprintf(
					"return value for method %s() is type [%s] but validator expects [%s]",
					name,
					retVal.String(),
					validator.Type(),
				),
			)
		}

		sv.getterValidators[name] = validator
	}

	// add friendly names for struct validators, if any
	if len(friendlyNames) > 0 {
		for _, names := range friendlyNames {
			for method, alias := range names {
				if _, ok := validators[method]; !ok {
					panic(fmt.Sprintf(`cannot set alias for method "%s()"; method does not exist`, method))
				}

				sv.getterAliases[method] = alias
			}
		}
	}

	return sv
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

	for _, check := range v.tests {
		if err := check(s); err != nil {
			return err
		}
	}

	// validate fields
	for field, val := range v.fieldValidators {
		fieldVal := sRef.FieldByName(field)
		if err := val.Validate(fieldVal.Interface()); err != nil {
			var printName string

			if alias, ok := v.fieldAliases[field]; ok {
				printName = alias
			} else {
				printName = field
			}

			return NewValidationError(fmt.Sprintf("%s: %s", printName, err.Error()))
		}
	}

	// validate getters
	for method, val := range v.getterValidators {
		methodVal := sRef.MethodByName(method)
		if err := val.Validate(methodVal.Interface()); err != nil {
			var printName string

			if alias, ok := v.getterAliases[method]; ok {
				printName = alias
			} else {
				printName = method
			}

			return NewValidationError(fmt.Sprintf("%s: %s", printName, err.Error()))
		}
	}
	return nil
}

// Is adds the provided function as a check against any values to be validated
func (v *StructValidator[T]) Is(fn structCheckFunc[T]) *StructValidator[T] {
	v.tests = append(v.tests, fn)
	return v
}
