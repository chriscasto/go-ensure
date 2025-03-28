package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

// TestAnyValidator_Construct checks to make sure construction fails with a panic
// when invalid inputs are provided
func TestAnyValidator_Construct(t *testing.T) {
	t.Run("panic if mismatched types", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Any(
			ensure.String(),
			ensure.Number[int](),
		)

		if err := bad.Validate(""); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}

	})
}

// TestAnyValidator_Type checks to make sure the AnyValidator returns the correct type
func TestAnyValidator_Type(t *testing.T) {
	testCases := map[string]struct {
		validator with.Validator
		t         string
	}{
		"string": {
			ensure.String(),
			"string",
		},
		"int": {
			ensure.Number[int](),
			"int",
		},
		"struct": {
			ensure.Struct[testStruct](),
			"ensure_test.testStruct",
		},
		"array of int": {
			ensure.Array[int](),
			"[]int",
		},
		"string pointer": {
			ensure.Pointer(ensure.String()),
			"*string",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			anyValid := ensure.Any(tc.validator)

			if anyValid.Type() != tc.t {
				t.Errorf(`unexpected type: expected "%s", got "%s"`, tc.t, anyValid.Type())
			}
		})
	}
}

// TestAnyValidator_WithError checks to make sure that the WithError method
// sets the return error correctly when validation fails
func TestAnyValidator_WithError(t *testing.T) {
	errMsg := "an error occurred"

	anyValid := ensure.Any(
		ensure.String().Equals("123"),
	).WithError(errMsg)

	if err := anyValid.Validate("abc"); err != nil {
		if err.Error() != errMsg {
			t.Errorf(`unexpected error: expected "%s", got "%s"`, errMsg, err)
		}
	}
}

func TestAnyValidator_Validate(t *testing.T) {
	testCases := map[string]struct {
		value    string
		willPass bool
	}{
		"match first":  {"foo", true},
		"match second": {"123", true},
		"match third":  {"validation", true},
		"match none":   {":(", false},
	}

	anyValid := ensure.Any(
		ensure.String().Equals("foo"),
		ensure.String().Matches(ensure.Numbers),
		ensure.String().IsLongerThan(5),
	)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := anyValid.Validate(tc.value)
			if err != nil && tc.willPass {
				t.Errorf(`expected no error, got "%s"`, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`expected error but got none`)
			}
		})
	}
}
