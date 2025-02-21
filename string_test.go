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

func TestStringValidator_Matches(t *testing.T) {
	t.Run("panic if bad regex", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.String().Matches("(")
		if err := bad.Validate(""); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}

	})

	testCases := strTestCases{
		"empty":      {"", false},
		"letters":    {"abc", true},
		"whitespace": {"abc xyz", false},
		"alphanum":   {"abc123", false},
		"upper":      {"ABC", false},
	}

	pattern := "^[a-z]+$"
	testCases.run(
		t,
		ensure.String().Matches(pattern),
		fmt.Sprintf(`Matches("%s")`, pattern),
	)

	constTests := map[string]strTestCases{
		ensure.Alpha: {
			"empty":      {"", false},
			"letters":    {"abc", true},
			"whitespace": {"abc xyz", false},
			"alphanum":   {"abc123", false},
			"upper":      {"ABC", true},
		},
		ensure.AlphaNum: {
			"empty":      {"", false},
			"letters":    {"abc", true},
			"whitespace": {"abc xyz", false},
			"alphanum":   {"abc123", true},
			"upper":      {"ABC", true},
		},
		ensure.Numbers: {
			"empty":    {"", false},
			"letters":  {"abc", false},
			"alphanum": {"abc123", false},
			"numbers":  {"123", true},
		},
		ensure.Decimal: {
			"empty":         {"", false},
			"letters":       {"abc", false},
			"alphanum":      {"abc123", false},
			"whole":         {"123", false},
			"no int":        {".123", true},
			"no fractional": {"123.", false},
			"dot":           {".", false},
		},
		ensure.Uuid4: {
			"empty":         {"", false},
			"good uuid":     {"cc3a331a-9899-46af-a2fd-dc4f73130d9c", true},
			"bad char uuid": {"xc3a331a-9899-46af-a2fd-dc4f73130d9c", false},
			"all bad uuid":  {"xy3h331u-9899-46li-o2pl-tr4h73130n9m", false},
		},
		ensure.Ipv4: {
			"empty":         {"", false},
			"good ip":       {"192.168.1.1", true},
			"bad num":       {"192.168.1.256", false},
			"missing octet": {"192.168.1.", false},
		},
		ensure.Email: {
			"empty":       {"", false},
			"simple":      {"test@example.com", true},
			"plus alias":  {"test+alias@example.com", true},
			"dot":         {"te.st@example.com", true},
			"upper":       {"TEST@EXAMPLE.COM", true},
			"empty name":  {"@example.com", false},
			"missing tld": {"test@example", false},
		},
		ensure.Md5: {
			"empty":         {"", false},
			"good hash":     {"d00413cdded7a5c5bc2e06079d63e562", true},
			"bad char hash": {"x00413cdded7a5c5bc2e06079d63e562", false},
			"short hash":    {"d00413cdded7a5c5bc2e06079d63e56", false},
		},
		ensure.Sha1: {
			"empty":         {"", false},
			"good hash":     {"03503cd82262dd65474b36365a5bc87dbdc12b25", true},
			"bad char hash": {"x3503cd82262dd65474b36365a5bc87dbdc12b25", false},
			"short hash":    {"03503cd82262dd65474b36365a5bc87dbdc12b2", false},
		},
		ensure.Sha256: {
			"empty": {"", false},
			"good hash": {
				"27ce252594811401f98751be0f91a3723651d2abdff7e33e99082d470b6c9184",
				true,
			},
			"bad char hash": {
				"x7ce252594811401f98751be0f91a3723651d2abdff7e33e99082d470b6c9184",
				false,
			},
			"short hash": {
				"27ce252594811401f98751be0f91a3723651d2abdff7e33e99082d470b6c918",
				false,
			},
		},
		ensure.Sha512: {
			"empty": {"", false},
			"good hash": {
				"9f45640d7fc448e8a7bc1398a030d19f7791b42e951a713e24d19b264c3fea4858818716586fd3470489095b2a7e81f47e79c951804c47ba111f7763d58c1de6",
				true,
			},
			"bad char hash": {
				"xf45640d7fc448e8a7bc1398a030d19f7791b42e951a713e24d19b264c3fea4858818716586fd3470489095b2a7e81f47e79c951804c47ba111f7763d58c1de6",
				false,
			},
			"short hash": {
				"9f45640d7fc448e8a7bc1398a030d19f7791b42e951a713e24d19b264c3fea4858818716586fd3470489095b2a7e81f47e79c951804c47ba111f7763d58c1de",
				false,
			},
		},
	}

	for constRegex, tc := range constTests {
		tc.run(
			t,
			ensure.String().Matches(constRegex),
			fmt.Sprintf(`Matches("%s")`, constRegex),
		)
	}
}

func TestStringValidator_StartsWith(t *testing.T) {
	testCases := strTestCases{
		"exact match":  {"foo", true},
		"substr match": {"food", true},
		"uppercase":    {"FOOD", false},
		"no match":     {"f", false},
	}

	prefix := "foo"
	testCases.run(
		t,
		ensure.String().StartsWith(prefix),
		fmt.Sprintf(`StartsWith("%s")`, prefix),
	)
}
