package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

// TestAnyValidator_IsValidator checks to make sure the AnyValidator implements the Validator interfaces
func TestAnyValidator_IsValidator(t *testing.T) {
	var _ with.UntypedValidator = ensure.Any[string]()
	var _ with.Validator[string] = ensure.Any[string]()
}

// TestAnyValidator_Type checks to make sure the AnyValidator returns the correct type
func TestAnyValidator_Type(t *testing.T) {
	testCases := map[string]struct {
		validator with.UntypedValidator
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

func TestAnyValidator_ValidateUntyped(t *testing.T) {
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
			err := anyValid.ValidateUntyped(tc.value)
			if err != nil && tc.willPass {
				t.Errorf(`expected no error, got "%s"`, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`expected error but got none`)
			}
		})
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

func TestAnyValidator_MultiError(t *testing.T) {
	intTestCases := multiErrTestCases[int]{
		"zero":  {0, 3}, // fails odd, greater than 1, equals 5
		"one":   {1, 2}, // fails greater than 1, equals 5
		"two":   {2, 2}, // fails odd, equals 5
		"three": {3, 1}, // fails equals 5
		"five":  {5, 0}, // fails none
		"six":   {6, 3}, // fails odd, less than 6, equals 5
	}

	intTestCases.run(t,
		ensure.Any[int](
			ensure.Number[int]().IsOdd().IsGreaterThan(1).IsLessThan(6).Equals(5),
		),
	)
}
