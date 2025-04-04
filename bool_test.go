package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

// TestBoolValidator_IsValidator checks to make sure the BoolValidator meets the Validator interface
func TestBoolValidator_IsValidator(t *testing.T) {
	// This should fail if BooleanValidator no longer meets the requirements for the Validator interface
	var _ with.Validator = ensure.Bool()
}

// TestBoolValidator_Type checks to make sure the BoolValidator returns the correct type
func TestBoolValidator_Type(t *testing.T) {
	bv := ensure.Bool()

	if bv.Type() != "bool" {
		t.Errorf(`unexpected type: expected "%s", got "%s"`, "bool", bv.Type())
	}
}

// TestBoolValidator_IsTrue checks to make sure the BoolValidator validates correctly when a true value is expected
func TestBoolValidator_IsTrue(t *testing.T) {
	bv := ensure.Bool().IsTrue()

	if err := bv.ValidateStrict(true); err != nil {
		t.Errorf(`expected no error, got "%s"`, err)
	}

	if err := bv.ValidateStrict(false); err == nil {
		t.Errorf(`expected error but got none`)
	}
}

// TestBoolValidator_IsFalse checks to make sure the BoolValidator validates correctly when a false value is expected
func TestBoolValidator_IsFalse(t *testing.T) {
	bv := ensure.Bool().IsFalse()

	if err := bv.ValidateStrict(false); err != nil {
		t.Errorf(`expected no error, got "%s"`, err)
	}

	if err := bv.ValidateStrict(true); err == nil {
		t.Errorf(`expected error but got none`)
	}
}

func TestBoolValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Bool())
}
