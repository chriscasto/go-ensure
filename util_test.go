package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"testing"
)

type testStruct struct {
	Str   string
	Int   int
	Float float64
}

type validatorTestCase struct {
	input    any
	willPass bool
}

type validatorTestCases map[string]*validatorTestCase

func (tcs *validatorTestCases) run(t *testing.T, v ensure.Validator) {
	for name, tc := range *tcs {
		t.Run(name, func(t *testing.T) {
			err := v.Validate(tc.input)
			if err != nil && tc.willPass {
				t.Errorf(`Validator[%s].Validate(%v) as {%s}; expected no error, got "%s"`, v.Type(), tc.input, name, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Validator[%s].Validate(%v) as {%s}; expected error but got none`, v.Type(), tc.input, name)
			}
		})
	}
}

func getDefaultValidatorTestCases(v ensure.Validator) validatorTestCases {
	testCases := validatorTestCases{
		"string":   {"a", false},
		"[]string": {[]string{"a", "b", "c"}, false},

		"int":   {1, false},
		"[]int": {[]int{1, 2, 3}, false},

		"float64":   {1.0, false},
		"[]float64": {[]float64{1.0, 2.0, 3.0}, false},

		"ensure_test.testStruct":   {testStruct{Str: "foo"}, false},
		"[]ensure_test.testStruct": {[]testStruct{{Str: "foo"}}, false},
	}

	// We expect any entry with a matching type to pass
	if testCases[v.Type()] != nil {
		testCases[v.Type()].willPass = true
	}

	return testCases
}

func runDefaultValidatorTestCases(t *testing.T, v ensure.Validator) {
	testCases := getDefaultValidatorTestCases(v)
	testCases.run(t, v)
}
