package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"strings"
	"testing"
)

// TestAnyValidator_IsValidator checks to make sure the AnyValidator implements the Validator interfaces
func TestAnyValidator_IsValidator(t *testing.T) {
	var _ with.UntypedValidator = ensure.Any[string](ensure.String())
	var _ with.Validator[string] = ensure.Any[string](ensure.String())
}

func TestAnyValidator_Construct(t *testing.T) {
	t.Run("not struct", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		constructBad := func(str string) error {
			bad := ensure.Any[string]()
			if err := bad.Validate(str); err != nil {
				return err
			}
			return nil
		}

		if err := constructBad("test"); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}
	})
}

// TestAnyValidator_Type checks to make sure the AnyValidator returns the correct type
func TestAnyValidator_Type(t *testing.T) {
	testCases := map[string]struct {
		validator with.UntypedValidator
		t         string
	}{
		"string": {
			ensure.Any[string](ensure.String()),
			"string",
		},
		"int": {
			ensure.Any[int](ensure.Number[int]()),
			"int",
		},
		"struct": {
			ensure.Any[testStruct](ensure.Struct[testStruct]()),
			"ensure_test.testStruct",
		},
		"array of int": {
			ensure.Any[[]int](ensure.Array[int]()),
			"[]int",
		},
		"string pointer": {
			ensure.Any[*string](ensure.Pointer[string](ensure.String())),
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

func TestAnyValidator_WithOptions_DefaultError(t *testing.T) {
	errMsg := "an error occurred"

	anyValid := ensure.Any[string](
		ensure.String().Equals("123"),
	).WithOptions(with.AnyOptionDefaultError(errMsg))

	if err := anyValid.Validate("abc"); err != nil {
		if err.Error() != errMsg {
			t.Errorf(`unexpected error: expected "%s", got "%s"`, errMsg, err)
		}
	}
}

func TestAnyValidator_WithOptions_PassThroughErrorsFrom(t *testing.T) {
	type anyTestStruct struct {
		Ignore bool
		Name   string
	}

	anyValid := ensure.Any[anyTestStruct](
		// If "Ignore" is true, consider that to be valid
		ensure.Struct[anyTestStruct]().HasFields(with.Validators{
			"Ignore": ensure.Bool().IsTrue(),
		}),
		// Otherwise, consider it valid if name is not empty and only contains alpha characters
		ensure.Struct[anyTestStruct]().HasFields(with.Validators{
			//"Ignore": ensure.Bool().IsFalse(),
			"Name": ensure.String().IsNotEmpty().Matches(ensure.Alpha),
		}),
	).WithOptions(
		// If validation fails, use the error(s) from the validator at index 1
		with.AnyOptionPassThroughErrorsFrom(1),
	)

	testCases := map[string]struct {
		val      anyTestStruct
		willPass bool
	}{
		"ignore": {
			val: anyTestStruct{
				Ignore: true,
			},
			willPass: true,
		},
		"valid name": {
			val: anyTestStruct{
				Name: "Alice",
			},
			willPass: true,
		},
		"invalid name": {
			val: anyTestStruct{
				Name: "Johnny 5",
			},
			willPass: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := anyValid.Validate(tc.val)

			if err != nil {
				if tc.willPass {
					t.Errorf(`expected no error, got "%s"`, err.Error())
				} else if !strings.Contains(err.Error(), "string") {
					t.Errorf(`expected error containing "string", got "%s"`, err.Error())
				}
			} else if err == nil && !tc.willPass {
				t.Errorf(`expected error, got none`)
			}
		})
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
	// These cases should only return a single error if we aren't passing errors through
	intTestCases1 := multiErrTestCases[int]{
		"zero":  {0, 1}, // fails odd, greater than 1, equals 5
		"one":   {1, 1}, // fails greater than 1, equals 5
		"two":   {2, 1}, // fails odd, equals 5
		"three": {3, 1}, // fails equals 5
		"five":  {5, 0}, // fails none
		"six":   {6, 1}, // fails odd, less than 6, equals 5
	}

	anyValid := ensure.Any[int](
		ensure.Number[int]().IsOdd().IsGreaterThan(1).IsLessThan(6).Equals(5),
	)

	intTestCases1.run(t, anyValid)

	// The same test cases will return multiple if we do pass through
	intTestCases2 := multiErrTestCases[int]{
		"zero":  {0, 3}, // fails odd, greater than 1, equals 5
		"one":   {1, 2}, // fails greater than 1, equals 5
		"two":   {2, 2}, // fails odd, equals 5
		"three": {3, 1}, // fails equals 5
		"five":  {5, 0}, // fails none
		"six":   {6, 3}, // fails odd, less than 6, equals 5
	}

	intTestCases2.run(t, anyValid.WithOptions(
		// If validation fails, use the error(s) from the first number validator
		with.AnyOptionPassThroughErrorsFrom(0),
	))
}
