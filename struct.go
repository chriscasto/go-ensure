package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

type structCheckFunc[T any] func(T) error

type methodValidator struct {
	method    *reflect.Method
	isPtr     bool
	validator with.Validator
}

type validField struct {
	name        string
	displayName string
	validator   with.Validator
}

type StructValidator[T any] struct {
	refVal  reflect.Value
	ptrType reflect.Type
	checks  []structCheckFunc[T]

	fields []*validField

	//fieldValidators with.Validators
	//fieldAliases    with.FriendlyNames

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
		refVal:  ref,
		ptrType: reflect.PointerTo(ref.Type()),
		fields:  []*validField{},
		//fieldValidators:  with.Validators{},
		//fieldAliases:     with.FriendlyNames{},
		getterValidators: with.Validators{},
		getterAliases:    with.FriendlyNames{},
	}
}

func (sv *StructValidator[T]) HasFields(validators with.Validators, friendlyNames ...with.FriendlyNames) *StructValidator[T] {
	ref := sv.refVal
	aliases := with.FriendlyNames{}

	// collect aliases for lookup during field processing
	if len(friendlyNames) > 0 {
		for _, names := range friendlyNames {
			for field, alias := range names {
				if _, ok := validators[field]; !ok {
					panic(fmt.Sprintf(`cannot set alias for field "%s"; field does not exist`, field))
				}

				aliases[field] = alias
			}
		}
	}

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

func (sv *StructValidator[T]) HasGetters(validators with.Validators, friendlyNames ...with.FriendlyNames) *StructValidator[T] {
	ref := sv.refVal

	// To get all the methods on a struct, we have to look at the pointer value
	ptr := reflect.PointerTo(ref.Type())

	//fmt.Println(ptr.NumMethod())
	//fmt.Println(ptr.String())
	//
	//for i := 0; i < ptr.NumMethod(); i++ {
	//	m := ptr.Method(i)
	//	fmt.Println(m.Name)
	//}

	//valMethods := map[string]*reflect.Method{}
	//ptrMethods := map[string]*reflect.Method{}

	// make sure that validator type matches field type
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

	ptrRef := reflect.PointerTo(sv.refVal.Type())

	// validate getters
	for method, val := range sv.getterValidators {
		methodVal, ok := ptrRef.MethodByName(method)

		// This shouldn't happen since we confirm method existence during HasGetter call
		if !ok {
			return NewValidationError(fmt.Sprintf("method %s does not exist", method))
		}

		// Call the getter method
		result := methodVal.Func.Call([]reflect.Value{reflect.ValueOf(&s)})
		retVal := result[0].Interface()

		if err := val.Validate(retVal); err != nil {
			var printName string

			if alias, ok := sv.getterAliases[method]; ok {
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
func (sv *StructValidator[T]) Is(fn structCheckFunc[T]) *StructValidator[T] {
	sv.checks = append(sv.checks, fn)
	return sv
}
