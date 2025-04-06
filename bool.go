package ensure

const boolType = "bool"

// BooleanValidator contains information and logic used to validate a boolean value
type BooleanValidator struct {
	expectTrue  bool
	expectFalse bool
}

// Bool returns an initialized BooleanValidator
func Bool() *BooleanValidator {
	return &BooleanValidator{
		expectTrue:  false,
		expectFalse: false,
	}
}

// Type will always return "bool"
func (bv *BooleanValidator) Type() string {
	return boolType
}

// IsTrue will cause the validator to return an error if the value is not true
func (bv *BooleanValidator) IsTrue() *BooleanValidator {
	bv.expectTrue = true
	return bv
}

// IsFalse will cause the validator to return an error if the value is not false
func (bv *BooleanValidator) IsFalse() *BooleanValidator {
	bv.expectFalse = true
	return bv
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a boolean
func (bv *BooleanValidator) ValidateUntyped(i interface{}) error {
	b, ok := i.(bool)

	if !ok {
		return NewTypeError("boolean expected")
	}

	return bv.Validate(b)
}

// Validate applies all checks against a boolean value and returns an error if any fail
func (bv *BooleanValidator) Validate(b bool) error {
	// There are really only two possibilities, so we can just check those
	// directly rather than using an array of functions

	if bv.expectTrue && b != true {
		return NewValidationError("expected true but got false")
	}

	if bv.expectFalse && b != false {
		return NewValidationError("expected false but got true")
	}

	return nil
}
