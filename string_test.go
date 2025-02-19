package ensure_test

import (
	"fmt"
	ensure "github.com/chriscasto/go-ensure"
	"testing"
)

type strTestCase struct {
	value    string
	willPass bool
}

type strTestCases map[string]strTestCase

func (tcs strTestCases) run(t *testing.T, sv *ensure.StringValidator, method string) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := sv.Validate(tc.value)
			if err != nil && tc.willPass {
				t.Errorf(`String().%s.Validate("%s"); expected no error, got "%s"`, method, tc.value, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`String().%s.Validate("%s"); expected error but got none`, method, tc.value)
			}
		})
	}
}

func TestStringValidator_Type(t *testing.T) {
	sv := ensure.String()
	sType := sv.Type()
	expected := "string"

	if sType != expected {
		t.Errorf("String.Type() = %s; want %s", sType, expected)
	}
}

func TestStringValidator_HasLength(t *testing.T) {
	testCases := strTestCases{
		"fewer letters": {"a", false},
		"same letters":  {"abc", true},
		"more letters":  {"wxyz", false},
	}

	strLen := 3
	testCases.run(
		t,
		ensure.String().HasLength(strLen),
		fmt.Sprintf("HasLength(%d)", strLen),
	)
}

func TestStringValidator_IsShorterThan(t *testing.T) {
	testCases := strTestCases{
		"fewer letters": {"a", true},
		"same letters":  {"abc", false},
		"more letters":  {"wxyz", false},
	}

	strLen := 3
	testCases.run(
		t,
		ensure.String().IsShorterThan(strLen),
		fmt.Sprintf("IsShorterThan(%d)", strLen),
	)
}

func TestStringValidator_IsLongerThan(t *testing.T) {
	testCases := strTestCases{
		"fewer letters": {"a", false},
		"same letters":  {"abc", false},
		"more letters":  {"wxyz", true},
	}

	strLen := 3
	testCases.run(
		t,
		ensure.String().IsLongerThan(strLen),
		fmt.Sprintf("IsLongerThan(%d)", strLen),
	)
}

func TestStringValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.String())
}

func TestShortenString(t *testing.T) {
	testCases := map[string]struct {
		input  string
		maxLen int
		want   string
	}{
		"short":      {"abc", 5, "abc"},
		"min maxLen": {"abc", 1, "abc"}, // min val for maxLen is 5
		"exact":      {"abcde", 5, "abcde"},
		"long":       {"abcdefghijklmnopqrstuvwxyz", 10, "abc...xyz"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			short := ensure.ShortenString(tc.input, tc.maxLen)

			if short != tc.want {
				t.Errorf(`got "%s"; want "%s" `, short, tc.want)
			}
		})
	}
}
