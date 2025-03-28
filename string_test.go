package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
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
			err := sv.ValidateString(tc.value)
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

func TestStringValidator_Length_Equals(t *testing.T) {
	testCases := strTestCases{
		"fewer letters": {"a", false},
		"same letters":  {"abc", true},
		"more letters":  {"wxyz", false},
	}

	strLen := 3
	testCases.run(
		t,
		ensure.String().HasLengthWhere(ensure.Length().Equals(strLen)),
		fmt.Sprintf("HasLengthWhere(Length().Equals(%d))", strLen),
	)
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

func TestStringValidator_Equals(t *testing.T) {
	testCases := strTestCases{
		"exact match":   {"foo", true},
		"partial match": {"food", false},
		"uppercase":     {"FOO", false},
	}

	str := "foo"
	testCases.run(
		t,
		ensure.String().Equals(str),
		fmt.Sprintf(`Equals("%s")`, str),
	)
}

func TestStringValidator_DoesNotEqual(t *testing.T) {
	testCases := strTestCases{
		"exact match":   {"foo", false},
		"partial match": {"food", true},
		"uppercase":     {"FOO", true},
	}

	str := "foo"
	testCases.run(
		t,
		ensure.String().DoesNotEqual(str),
		fmt.Sprintf(`DoesNotEqual("%s")`, str),
	)
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

func TestStringValidator_DoesNotStartWith(t *testing.T) {
	testCases := strTestCases{
		"exact match":  {"foo", false},
		"substr match": {"food", false},
		"uppercase":    {"FOOD", true},
		"no match":     {"f", true},
	}

	prefix := "foo"
	testCases.run(
		t,
		ensure.String().DoesNotStartWith(prefix),
		fmt.Sprintf(`DoesNotStartWith("%s")`, prefix),
	)
}

func TestStringValidator_EndsWith(t *testing.T) {
	testCases := strTestCases{
		"exact match":  {"bar", true},
		"substr match": {"crowbar", true},
		"uppercase":    {"CROWBAR", false},
		"mixed case":   {"CROWBAr", false},
		"no match":     {"rab", false},
	}

	suffix := "bar"
	testCases.run(
		t,
		ensure.String().EndsWith(suffix),
		fmt.Sprintf(`EndsWith("%s")`, suffix),
	)
}

func TestStringValidator_DoesNotEndWith(t *testing.T) {
	testCases := strTestCases{
		"exact match":  {"bar", false},
		"substr match": {"crowbar", false},
		"uppercase":    {"CROWBAR", true},
		"mixed case":   {"CROWBAr", true},
		"no match":     {"rab", true},
	}

	suffix := "bar"
	testCases.run(
		t,
		ensure.String().DoesNotEndWith(suffix),
		fmt.Sprintf(`DoesNotEndWith("%s")`, suffix),
	)
}

func TestStringValidator_Contains(t *testing.T) {
	testCases := strTestCases{
		"exact match": {"boo", true},
		"prefix":      {"book", true},
		"suffix":      {"taboo", true},
		"uppercase":   {"BOO", false},
		"mixed case":  {"TaBoO", false},
		"no match":    {"bu", false},
	}

	substr := "boo"
	testCases.run(
		t,
		ensure.String().Contains(substr),
		fmt.Sprintf(`Contains("%s")`, substr),
	)
}

func TestStringValidator_DoesNotContain(t *testing.T) {
	testCases := strTestCases{
		"exact match": {"boo", false},
		"prefix":      {"book", false},
		"suffix":      {"taboo", false},
		"uppercase":   {"BOO", true},
		"mixed case":  {"TaBoO", true},
		"no match":    {"bu", true},
	}

	substr := "boo"
	testCases.run(
		t,
		ensure.String().DoesNotContain(substr),
		fmt.Sprintf(`DoesNotContain("%s")`, substr),
	)
}

func TestStringValidator_IsEmpty(t *testing.T) {
	testCases := strTestCases{
		"empty":      {"", true},
		"whitespace": {"  ", false},
		"not empty":  {"abc", false},
	}

	testCases.run(
		t,
		ensure.String().IsEmpty(),
		"IsEmpty()",
	)
}

func TestStringValidator_IsNotEmpty(t *testing.T) {
	testCases := strTestCases{
		"empty":      {"", false},
		"whitespace": {"  ", true},
		"not empty":  {"abc", true},
	}

	testCases.run(
		t,
		ensure.String().IsNotEmpty(),
		"IsNotEmpty()",
	)
}

func TestStringValidator_IsOneOf(t *testing.T) {
	testCases := strTestCases{
		"in set":     {"one", true},
		"not in set": {"two", false},
		"upper":      {"ONE", false},
	}

	permitted := []string{
		"one",
		"three",
	}
	testCases.run(
		t,
		ensure.String().IsOneOf(permitted),
		fmt.Sprintf(`IsOneOf(%v)`, permitted),
	)
}

func TestStringValidator_IsNotOneOf(t *testing.T) {
	testCases := strTestCases{
		"in set":     {"one", false},
		"not in set": {"two", true},
		"upper":      {"ONE", true},
	}

	forbidden := []string{
		"one",
		"three",
	}
	testCases.run(
		t,
		ensure.String().IsNotOneOf(forbidden),
		fmt.Sprintf(`IsNotOneOf(%v)`, forbidden),
	)
}
