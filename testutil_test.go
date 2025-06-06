package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"strings"
	"testing"
)

// testStruct is used to test struct validation
type testStruct struct {
	Str   string
	Int   int
	Float float64
}

// GetStr is used to test getter validation on a string type
func (ts testStruct) GetStr() string {
	return ts.Str
}

// GetInt is used to test getter validation on an int type
// Note that the use of a ptr receiver here is intentional to ensure both receiver types work as expected
func (ts *testStruct) GetInt() int {
	return ts.Int
}

// GetFloat is used to test getter validation on a float type
func (ts testStruct) GetFloat() float64 {
	return ts.Float
}

// GetStrWithArg is used to test that getter validation fails if method has an arg
func (ts testStruct) GetStrWithArg(upper bool) string {
	if upper {
		return strings.ToUpper(ts.Str)
	}
	return ts.Str
}

// GetStrWithError is used to test that getter validation fails if method returns multiple values
func (ts testStruct) GetStrWithError() (string, error) {
	return ts.Str, nil
}

type validatorTestCase struct {
	input    any
	willPass bool
}

type validatorTestCases map[string]*validatorTestCase

func (tcs *validatorTestCases) run(t *testing.T, v with.UntypedValidator) {
	for name, tc := range *tcs {
		t.Run(name, func(t *testing.T) {
			err := v.ValidateUntyped(tc.input)
			if err != nil && tc.willPass {
				t.Errorf(`Validator[%s].Validate(%v) as {%s}; expected no error, got "%s"`, v.Type(), tc.input, name, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Validator[%s].Validate(%v) as {%s}; expected error but got none`, v.Type(), tc.input, name)
			}
		})
	}
}

// getDefaultValidatorTestCases generates a set of test cases to confirm that
// validators only accept values of the correct type
func getDefaultValidatorTestCases(v with.UntypedValidator) validatorTestCases {
	testCases := validatorTestCases{
		"bool":   {true, false},
		"[]bool": {[]bool{true, false}, false},

		"string":   {"a", false},
		"[]string": {[]string{"a", "b", "c"}, false},

		"int":   {1, false},
		"[]int": {[]int{1, 2, 3}, false},

		"float64":   {1.0, false},
		"[]float64": {[]float64{1.0, 2.0, 3.0}, false},

		"ensure_test.testStruct":   {testStruct{Str: "foo"}, false},
		"[]ensure_test.testStruct": {[]testStruct{{Str: "foo"}}, false},

		"map[string]int":        {map[string]int{"a": 1, "b": 2}, false},
		"map[string][]int":      {map[string][]int{"a": {1, 2, 3}}, false},
		"map[string]testStruct": {map[string][]testStruct{"a": {}}, false},
	}

	// We expect any entry with a matching type to pass
	if testCases[v.Type()] != nil {
		testCases[v.Type()].willPass = true
	}

	return testCases
}

func runDefaultValidatorTestCases(t *testing.T, v with.UntypedValidator) {
	testCases := getDefaultValidatorTestCases(v)
	testCases.run(t, v)
}

type multiErrTestCases[T any] map[string]struct {
	value            T
	expectedErrCount int
}

// run evaluates all the test cases that expect multiple errors
func (tcs multiErrTestCases[T]) run(t *testing.T, sv with.Validator[T]) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := sv.Validate(tc.value, with.Options(
				with.OptionCollectAllErrors(),
			))

			// we expect validation to return a total number of errors equal to tc.expectedErrCount

			if err != nil {
				// we got an error, but didn't expect any
				if tc.expectedErrCount == 0 {
					t.Errorf(`expected no error, got "%s"`, err)
				} else {
					// we expected some error and got one
					// now we see if we can unwrap that into multiple errors
					vErrs := ensure.ErrorAsValidationErrors(err)

					if vErrs == nil {
						t.Errorf(`expected validation errors but got unexpected error type`)
					} else {
						errList := vErrs.ValidationErrors()

						if len(errList) != tc.expectedErrCount {
							t.Errorf(`expected %d errors but got %d`, tc.expectedErrCount, len(errList))
						}
					}
				}
			} else {
				// no error, but we expected some
				if tc.expectedErrCount != 0 {
					t.Errorf(`expected %d errors, got none`, tc.expectedErrCount)
				}
			}
		})
	}
}
