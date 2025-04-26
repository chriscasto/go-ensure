package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"regexp"
	"strings"
)

//goland:noinspection GoCommentStart
const (
	// Basic patterns
	Alpha    = `(?i)^[a-z]+$`
	AlphaNum = `(?i)^[a-z0-9]+$`
	Numbers  = `^\d+$`
	Decimal  = `^\d*\.\d+$`

	// Uuid
	Uuid4 = `(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

	// Internet
	Ipv4  = `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	Email = `(?i)(?:[a-z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`

	// Hashes
	Md5    = `^[0-9a-f]{32}$`
	Sha1   = `^[0-9a-f]{40}$`
	Sha256 = `^[0-9a-f]{64}$`
	Sha512 = `^[0-9a-f]{128}$`
)

// StringValidator contains information and logic used to validate a string
type StringValidator struct {
	checks *lenChecks[string, string, string]
}

// String returns an initialized StringValidator
func String() *StringValidator {
	return &StringValidator{
		checks: newLenChecks[string, string, string](),
	}
}

// Type returns the string "string"
func (v *StringValidator) Type() string {
	return "string"
}

// HasLengthWhere adds a NumberValidator for validating the length of the string
func (v *StringValidator) HasLengthWhere(nv *NumberValidator[int]) *StringValidator {
	v.checks.AddHasLengthWhere(nv)
	return v
}

// ValidateUntyped accepts an arbitrary input type and validates it if it's a match for the expected type
func (v *StringValidator) ValidateUntyped(value any, options ...*with.ValidationOptions) error {
	str, ok := value.(string)

	if !ok {
		return NewTypeError("string expected")
	}

	return v.Validate(str, options...)
}

// Validate applies all checks against a string value and returns an error if any fail
func (v *StringValidator) Validate(str string, options ...*with.ValidationOptions) error {
	return v.checks.Evaluate(str, getValidationOptions(options))
}

// Equals adds a validation check that returns an error if the target string
// is not identical to the specified string
func (v *StringValidator) Equals(same string) *StringValidator {
	return v.Is(func(str string) error {
		if str != same {
			return errors.New(
				fmt.Sprintf(`string must equal "%s"`, same),
			)
		}
		return nil
	})
}

// DoesNotEqual adds a validation check that returns an error if the target string
// is identical to the specified string
func (v *StringValidator) DoesNotEqual(diff string) *StringValidator {
	return v.Is(func(str string) error {
		if str == diff {
			return errors.New(
				fmt.Sprintf(`string must not equal "%s"`, diff),
			)
		}
		return nil
	})
}

// StartsWith adds a validation check that returns an error if the target string
// does not start with the specified substring
func (v *StringValidator) StartsWith(prefix string) *StringValidator {
	return v.Is(func(str string) error {
		if !strings.HasPrefix(str, prefix) {
			return errors.New(
				fmt.Sprintf(`string must start with "%s"`, prefix),
			)
		}
		return nil
	})
}

// DoesNotStartWith adds a validation check that returns an error if the target string
// starts with the specified substring
func (v *StringValidator) DoesNotStartWith(prefix string) *StringValidator {
	return v.Is(func(str string) error {
		if strings.HasPrefix(str, prefix) {
			return errors.New(
				fmt.Sprintf(`string must not start with "%s"`, prefix),
			)
		}
		return nil
	})
}

// EndsWith adds a validation check that returns an error if the target string
// does not end with the specified substring
func (v *StringValidator) EndsWith(suffix string) *StringValidator {
	return v.Is(func(str string) error {
		if !strings.HasSuffix(str, suffix) {
			return errors.New(
				fmt.Sprintf(`string must end with "%s"`, suffix),
			)
		}
		return nil
	})
}

// DoesNotEndWith adds a validation check that returns an error if the target string
// ends with the specified substring
func (v *StringValidator) DoesNotEndWith(suffix string) *StringValidator {
	return v.Is(func(str string) error {
		if strings.HasSuffix(str, suffix) {
			return errors.New(
				fmt.Sprintf(`string must not end with "%s"`, suffix),
			)
		}
		return nil
	})
}

// Contains adds a validation check that returns an error if the target string
// does not contain the specified substring
func (v *StringValidator) Contains(substr string) *StringValidator {
	return v.Is(func(str string) error {
		if !strings.Contains(str, substr) {
			return errors.New(
				fmt.Sprintf(`string must contain "%s"`, substr),
			)
		}
		return nil
	})
}

// DoesNotContain adds a validation check that returns an error if the target string
// contains the specified substring
func (v *StringValidator) DoesNotContain(substr string) *StringValidator {
	return v.Is(func(str string) error {
		if strings.Contains(str, substr) {
			return errors.New(
				fmt.Sprintf(`string must not contain "%s"`, substr),
			)
		}
		return nil
	})
}

// IsEmpty adds a validation check that returns an error if the target string is not empty
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(0))
func (v *StringValidator) IsEmpty() *StringValidator {
	v.checks.AddIsEmpty()
	return v
}

// IsNotEmpty adds a validation check that returns an error if the target string is empty
// This is a convenience function that is equivalent to HasLengthWhere(Length().DoesNotEqual(0))
func (v *StringValidator) IsNotEmpty() *StringValidator {
	v.checks.AddIsNotEmpty()
	return v
}

// IsOneOf adds a validation check that returns an error if the target string
// is not in the specified set
func (v *StringValidator) IsOneOf(values []string) *StringValidator {
	// convert list to map for O(1) lookups
	lookup := map[string]bool{}

	for _, str := range values {
		lookup[str] = true
	}

	return v.Is(func(str string) error {
		if _, ok := lookup[str]; !ok {
			return errors.New(`string must be one of the permitted values`)
		}
		return nil
	})
}

// IsNotOneOf adds a validation check that returns an error if the target string
// is in the specified set
func (v *StringValidator) IsNotOneOf(values []string) *StringValidator {
	// convert list to map for O(1) lookups
	lookup := map[string]bool{}

	for _, str := range values {
		lookup[str] = true
	}

	return v.Is(func(str string) error {
		if _, ok := lookup[str]; ok {
			return errors.New(`string must not be one of the prohibited values`)
		}
		return nil
	})
}

// IsLongerThan adds a validation check that returns an error if the target
// string length is less than or equal to the specified value
func (v *StringValidator) IsLongerThan(l int) *StringValidator {
	v.checks.AddIsLongerThan(l)
	return v
}

// IsShorterThan adds a validation check that returns an error if the target
// string length is greater than or equal to the specified value
func (v *StringValidator) IsShorterThan(l int) *StringValidator {
	v.checks.AddIsShorterThan(l)
	return v
}

// HasLength adds a check that returns an error if the length of the string does not equal the provided value
// This is a convenience function that is equivalent to HasLengthWhere(Length().Equals(l))
func (v *StringValidator) HasLength(l int) *StringValidator {
	v.checks.AddHasLength(l)
	return v
}

// Matches adds a validation check that returns an error if the target
// string does not match the specified pattern
func (v *StringValidator) Matches(pattern string) *StringValidator {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("could not compile regex: %s", err))
	}

	return v.Is(func(str string) error {
		if !r.MatchString(str) {
			return errors.New(`string does not match expected pattern`)
		}
		return nil
	})
}

// Is adds the provided function as a check against any values to be validated
func (v *StringValidator) Is(fn func(string) error) *StringValidator {
	v.checks.Append(func(val string, _ *with.ValidationOptions) error {
		return fn(val)
	})
	return v
}
