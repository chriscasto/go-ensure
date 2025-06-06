package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"reflect"
)

// validMethod contains information about a method that needs to be called during validation
type validMethod struct {
	ref         *reflect.Method
	displayName string
	hasValRcvr  bool
	validator   with.UntypedValidator
}

// validField contains information about a field that needs to be accessed during validation
type validField struct {
	name        string
	displayName string
	validator   with.UntypedValidator
}

// StructValidator contains information and logic used to validate a struct of type T
type StructValidator[T any] struct {
	refVal  reflect.Value
	checks  *valChecks[T]
	fields  []*validField
	getters []*validMethod
}

// Struct constructs a StructValidator instance of type T and returns a pointer to it
func Struct[T any]() *StructValidator[T] {
	// Create an empty instance of the struct
	var zero T

	ref := reflect.ValueOf(zero)

	if ref.Kind() != reflect.Struct {
		panic("StructValidator expects a struct type")
	}

	return &StructValidator[T]{
		refVal:  ref,
		fields:  []*validField{},
		getters: []*validMethod{},
		checks:  newValChecks[T](),
	}
}

// Type returns a string with the name of the struct this validator expects
func (sv *StructValidator[T]) Type() string {
	return sv.refVal.Type().String()
}

// Is adds the provided function as a check against any values to be validated
func (sv *StructValidator[T]) Is(fn func(T) error) *StructValidator[T] {
	sv.checks.Append(func(val T, _ *with.ValidationOptions) error {
		return fn(val)
	})
	return sv
}

// Has adds the provided function as a check against any values to be validated
// Has is an alias for Is
func (sv *StructValidator[T]) Has(fn func(T) error) *StructValidator[T] {
	return sv.Is(fn)
}

// HasFields accepts a map of named fields and their validators to evaluate against a struct during validation
// It also accepts an optional map of field names to display names to use when printing error messages
func (sv *StructValidator[T]) HasFields(validators with.Validators, displayNames ...with.DisplayNames) *StructValidator[T] {
	ref := sv.refVal
	aliases := with.DisplayNames{}

	// Collect aliases for lookup during field processing
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

// HasGetters accepts a map of named getter methods and their validators to evaluate against a struct during validation
// It also accepts an optional map of method names to display names to use when printing error messages
func (sv *StructValidator[T]) HasGetters(validators with.Validators, displayNames ...with.DisplayNames) *StructValidator[T] {
	refType := sv.refVal.Type()

	// To get all the methods on a struct, we have to look at the pointer value
	ptr := reflect.PointerTo(refType)

	aliases := with.DisplayNames{}

	// Collect aliases for lookup during method processing
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

	// Collect all methods with value receivers
	// We use this as a simple check to know whether we should pass a value or a pointer at evaluation time
	valMethods := map[string]reflect.Method{}

	for i := range refType.NumMethod() {
		method := refType.Method(i)
		valMethods[method.Name] = method
	}

	for name, validator := range validators {
		// First check to see if the method has a value type receiver
		method, ok := valMethods[name]
		hasValRcvr := true

		// If we don't find it there, check the ptr version
		if !ok {
			hasValRcvr = false
			method, ok = ptr.MethodByName(name)
		}

		if !ok {
			panic(
				fmt.Sprintf(
					"method %s() does not exist in struct %s",
					name,
					refType.String(),
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
			hasValRcvr:  hasValRcvr,
			validator:   validator,
		})
	}

	return sv
}

// validateStruct is a helper method that does the actual validation used by Validate and ValidateStrict
func (sv *StructValidator[T]) validateStruct(sRef reflect.Value, s T, options ...*with.ValidationOptions) error {
	vErrs := newValidationErrors()
	vOpts := getValidationOptions(options)

	if err := sv.checks.Evaluate(s, vOpts); err != nil {
		if vOpts.CollectAllErrors() {
			vErrs.Append(err)
		} else {
			return err
		}
	}

	// Validate fields
	for _, field := range sv.fields {
		fieldVal := sRef.FieldByName(field.name)
		if err := field.validator.ValidateUntyped(fieldVal.Interface(), vOpts); err != nil {
			vErrs.Append(fmt.Errorf("%s: %s", field.displayName, err.Error()))

			if !vOpts.CollectAllErrors() {
				return vErrs
			}
		}
	}

	// Validate getters
	for _, method := range sv.getters {
		var receiver reflect.Value

		if method.hasValRcvr {
			receiver = reflect.ValueOf(s)
		} else {
			receiver = reflect.ValueOf(&s)
		}

		// Call the getter method
		result := method.ref.Func.Call([]reflect.Value{receiver})
		retVal := result[0].Interface()

		if err := method.validator.ValidateUntyped(retVal, vOpts); err != nil {
			vErrs.Append(fmt.Errorf("%s: %s", method.displayName, err.Error()))

			if !vOpts.CollectAllErrors() {
				return vErrs
			}
		}
	}

	if vErrs.HasErrors() {
		return vErrs
	}

	return nil
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a match for the expected type
func (sv *StructValidator[T]) ValidateUntyped(value any, options ...*with.ValidationOptions) error {
	sRef := reflect.ValueOf(value)
	sRefType := sRef.Type()

	if !sRef.IsValid() || sRefType != sv.refVal.Type() {
		return newTypeErrorFromTypes(sv.refVal.Type().String(), sRefType.String())
	}

	return sv.validateStruct(sRef, value.(T), options...)
}

// Validate applies all checks against a struct of the expected type and returns an error if any fail
func (sv *StructValidator[T]) Validate(s T, options ...*with.ValidationOptions) error {
	sRef := reflect.ValueOf(s)
	return sv.validateStruct(sRef, s, options...)
}
