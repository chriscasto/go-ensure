package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

type structCheckFunc[T any] func(T) error

type validMethod struct {
	ref         *reflect.Method
	displayName string
	isPtr       bool
	validator   with.Validator
}

type validField struct {
	name        string
	displayName string
	validator   with.Validator
}

type StructValidator[T any] struct {
	refVal  reflect.Value
	checks  []structCheckFunc[T]
	fields  []*validField
	getters []*validMethod
}

func Struct[T any]() *StructValidator[T] {
	// Create an empty instance of the struct
	var zero T

	ref := reflect.ValueOf(zero)

	if ref.Kind() != reflect.Struct {
		panic("type T must be a struct")
	}

	return &StructValidator[T]{
		refVal:  ref,
		fields:  []*validField{},
		getters: []*validMethod{},
	}
}

func (sv *StructValidator[T]) HasFields(validators with.Validators, displayNames ...with.DisplayNames) *StructValidator[T] {
	ref := sv.refVal
	aliases := with.DisplayNames{}

	// collect aliases for lookup during field processing
	if len(displayNames) > 0 {
		for _, names := range displayNames {
			for field, alias := range names {
				if _, ok := validators[field]; !ok {
					panic(fmt.Sprintf(`cannot set display name for field "%s"; field is not in list of validators`, field))
				}

				aliases[field] = alias
			}
		}
	}

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

		displayName, ok := aliases[name]

		if !ok {
			displayName = name
		}

		sv.fields = append(sv.fields, &validField{
			name:        name,
			validator:   validator,
			displayName: displayName,
		})
	}

	return sv
}

func (sv *StructValidator[T]) HasGetters(validators with.Validators, displayNames ...with.DisplayNames) *StructValidator[T] {
	ref := sv.refVal

	// To get all the methods on a struct, we have to look at the pointer value
	ptr := reflect.PointerTo(ref.Type())

	aliases := with.DisplayNames{}

	// collect aliases for lookup during method processing
	if len(displayNames) > 0 {
		for _, names := range displayNames {
			for method, alias := range names {
				if _, ok := validators[method]; !ok {
					panic(fmt.Sprintf(`cannot set display name for method "%s()"; method is not in list of validators`, method))
				}

				aliases[method] = alias
			}
		}
	}

	for name, validator := range validators {
		method, ok := ptr.MethodByName(name)

		if !ok {
			panic(
				fmt.Sprintf(
					"method %s() does not exist in struct %s",
					name,
					ref.Type().String(),
				),
			)
		}

		mType := method.Type

		// Getters can only have a single arg (the receiver)
		if mType.NumIn() != 1 {
			panic(
				fmt.Sprintf(
					"method %s() has %d args but validator expects one (receiver)",
					name,
					mType.NumIn(),
				),
			)
		}

		// Getters must only return a single value
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

		// The getter must return a value that matches the validator
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

		displayName, ok := aliases[name]

		if !ok {
			displayName = name
		}

		sv.getters = append(sv.getters, &validMethod{
			displayName: displayName,
			ref:         &method,
			validator:   validator,
		})
	}

	return sv
}

func (sv *StructValidator[T]) Type() string {
	return sv.refVal.Type().String()
}

func (sv *StructValidator[T]) Validate(s interface{}) error {
	sRef := reflect.ValueOf(s)
	sRefType := sRef.Type()

	if !sRef.IsValid() || sRefType != sv.refVal.Type() {
		return newTypeErrorFromTypes(sv.refVal.Type().String(), sRefType.String())
	}

	return sv.ValidateStruct(s.(T))
}

func (sv *StructValidator[T]) ValidateStruct(s T) error {
	sRef := reflect.ValueOf(s)

	for _, check := range sv.checks {
		if err := check(s); err != nil {
			return err
		}
	}

	// validate fields
	for _, field := range sv.fields {
		fieldVal := sRef.FieldByName(field.name)
		if err := field.validator.Validate(fieldVal.Interface()); err != nil {
			return NewValidationError(fmt.Sprintf("%s: %s", field.displayName, err.Error()))
		}
	}

	// validate getters
	for _, method := range sv.getters {
		// Call the getter method
		result := method.ref.Func.Call([]reflect.Value{reflect.ValueOf(&s)})
		retVal := result[0].Interface()

		if err := method.validator.Validate(retVal); err != nil {
			return NewValidationError(fmt.Sprintf("%s: %s", method.displayName, err.Error()))
		}
	}

	return nil
}

// Is adds the provided function as a check against any values to be validated
func (sv *StructValidator[T]) Is(fn structCheckFunc[T]) *StructValidator[T] {
	sv.checks = append(sv.checks, fn)
	return sv
}
