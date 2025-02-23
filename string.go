package ensure

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

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

type StringValidator struct {
	zeroVal string
	tests   []func(string) error
}

func String() *StringValidator {
	return &StringValidator{}
}

func (v *StringValidator) Type() string {
	return "string"
}

func (v *StringValidator) Validate(i interface{}) error {
	str, ok := i.(string)

	if !ok {
		return &TypeError{"string expected"}
	}

	for _, fn := range v.tests {
		if err := fn(str); err != nil {
			return NewValidationError(err.Error())
		}
	}

	return nil
}

// StartsWith adds a validation check that returns an error if the target string
// does not start with the specified substring
func (v *StringValidator) StartsWith(prefix string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if !strings.HasPrefix(str, prefix) {
			return errors.New(
				fmt.Sprintf(`string must start with "%s"`, prefix),
			)
		}
		return nil
	})
	return v
}

// DoesNotStartWith adds a validation check that returns an error if the target string
// starts with the specified substring
func (v *StringValidator) DoesNotStartWith(prefix string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if strings.HasPrefix(str, prefix) {
			return errors.New(
				fmt.Sprintf(`string must not start with "%s"`, prefix),
			)
		}
		return nil
	})
	return v
}

// EndsWith adds a validation check that returns an error if the target string
// does not end with the specified substring
func (v *StringValidator) EndsWith(suffix string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if !strings.HasSuffix(str, suffix) {
			return errors.New(
				fmt.Sprintf(`string must end with "%s"`, suffix),
			)
		}
		return nil
	})
	return v
}

// DoesNotEndWith adds a validation check that returns an error if the target string
// ends with the specified substring
func (v *StringValidator) DoesNotEndWith(suffix string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if strings.HasSuffix(str, suffix) {
			return errors.New(
				fmt.Sprintf(`string must not end with "%s"`, suffix),
			)
		}
		return nil
	})
	return v
}

// Contains adds a validation check that returns an error if the target string
// does not contain the specified substring
func (v *StringValidator) Contains(substr string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if !strings.Contains(str, substr) {
			return errors.New(
				fmt.Sprintf(`string must contain "%s"`, substr),
			)
		}
		return nil
	})
	return v
}

// DoesNotContain adds a validation check that returns an error if the target string
// contains the specified substring
func (v *StringValidator) DoesNotContain(substr string) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if strings.Contains(str, substr) {
			return errors.New(
				fmt.Sprintf(`string must not contain "%s"`, substr),
			)
		}
		return nil
	})
	return v
}

// IsEmpty adds a validation check that returns an error if the target string is not empty
func (v *StringValidator) IsEmpty() *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) != 0 {
			return errors.New(
				fmt.Sprintf(`string must be empty`),
			)
		}
		return nil
	})
	return v
}

// IsNotEmpty adds a validation check that returns an error if the target string is empty
func (v *StringValidator) IsNotEmpty() *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) == 0 {
			return errors.New(
				fmt.Sprintf(`string must not be empty`),
			)
		}
		return nil
	})
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

	v.tests = append(v.tests, func(str string) error {
		if _, ok := lookup[str]; !ok {
			return errors.New(
				fmt.Sprintf(`string must be one of the permitted values`),
			)
		}
		return nil
	})
	return v
}

// IsNotOneOf adds a validation check that returns an error if the target string
// is in the specified set
func (v *StringValidator) IsNotOneOf(values []string) *StringValidator {
	// convert list to map for O(1) lookups
	lookup := map[string]bool{}

	for _, str := range values {
		lookup[str] = true
	}

	v.tests = append(v.tests, func(str string) error {
		if _, ok := lookup[str]; ok {
			return errors.New(
				fmt.Sprintf(`string must not be one of the prohibited values`),
			)
		}
		return nil
	})
	return v
}

// IsLongerThan adds a validation check that returns an error if the target
// string length is less than or equal to the specified value
func (v *StringValidator) IsLongerThan(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) <= l {
			return errors.New(
				fmt.Sprintf(`string length must be greater than %d`, l),
			)
		}
		return nil
	})
	return v
}

// IsLongerThanOrEqualTo adds a validation check that returns an error if the target
// string length is less than the specified value
func (v *StringValidator) IsLongerThanOrEqualTo(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) < l {
			return errors.New(
				fmt.Sprintf(`string length must be greater than or equal to %d`, l),
			)
		}
		return nil
	})
	return v
}

// IsShorterThan adds a validation check that returns an error if the target
// string length is greater than or equal to the specified value
func (v *StringValidator) IsShorterThan(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) >= l {
			return errors.New(
				fmt.Sprintf(`string length must be less than %d`, l),
			)
		}
		return nil
	})
	return v
}

// IsShorterThanOrEqualTo adds a validation check that returns an error if the target
// string length is greater than the specified value
func (v *StringValidator) IsShorterThanOrEqualTo(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) > l {
			return errors.New(
				fmt.Sprintf(`string length must be less than or equal to %d`, l),
			)
		}
		return nil
	})
	return v
}

// HasLength adds a validation check that returns an error if the target
// string length does not equal the specified value
func (v *StringValidator) HasLength(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) != l {
			return errors.New(
				fmt.Sprintf(`string must have a length of exactly %d`, l),
			)
		}
		return nil
	})
	return v
}

// Matches adds a validation check that returns an error if the target
// string does not match the specified pattern
func (v *StringValidator) Matches(pattern string) *StringValidator {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	v.tests = append(v.tests, func(str string) error {
		if !r.MatchString(str) {
			return errors.New(
				fmt.Sprintf(
					`string does not match expected pattern`,
				),
			)
		}
		return nil
	})
	return v
}
