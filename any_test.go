package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

// TestAnyValidator_Type checks to make sure the AnyValidator returns the correct type
func TestAnyValidator_Type(t *testing.T) {
	testCases := map[string]struct {
		validator with.Validator
		t         string
	}{
		"string": {
			ensure.Any[string](),
			"string",
		},
		"int": {
			ensure.Any[int](),
			"int",
		},
		"struct": {
			ensure.Any[testStruct](),
			"ensure_test.testStruct",
		},
		"array of int": {
			ensure.Any[[]int](),
			"[]int",
		},
		"string pointer": {
			ensure.Any[*string](),
			"*string",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.validator.Type() != tc.t {
				t.Errorf(`unexpected type: expected "%s", got "%s"`, tc.t, tc.validator.Type())
			}
		})
	}
}

// TestAnyValidator_WithError checks to make sure that the WithError method
// sets the return error correctly when validation fails
func TestAnyValidator_WithError(t *testing.T) {
	errMsg := "an error occurred"

	anyValid := ensure.Any[string](
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

	anyValid := ensure.Any[string](
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
