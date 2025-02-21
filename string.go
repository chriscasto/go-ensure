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

// define the max length for string contents in an error message
const defaultMaxPrintLength = 30

// shortens a string of arbitrary length to a reasonable value for logging
func shortenString(s string, maxLen int) string {
	if maxLen < 5 {
		maxLen = 5
	}

	if len(s) <= maxLen {
		return s
	}

	ellipsis := "..."

	half := (maxLen - len(ellipsis)) / 2

	return s[:half] + ellipsis + s[len(s)-half:]
}

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
		return fmt.Errorf("string expected")
	}

	for _, fn := range v.tests {
		err := fn(str)
		if err != nil {
			return err
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
				fmt.Sprintf(
					`string "%s" does not contain prefix "%s"`,
					shortenString(str, defaultMaxPrintLength),
					prefix),
			)
		}

		return nil
	})

	return v
}

func (v *StringValidator) IsLongerThan(l int) *StringValidator {

	v.tests = append(v.tests, func(str string) error {
		if len(str) <= l {
			return errors.New(
				fmt.Sprintf(
					`string "%s" is shorter than min length (%d)`,
					shortenString(str, defaultMaxPrintLength),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *StringValidator) IsShorterThan(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) >= l {
			return errors.New(
				fmt.Sprintf(
					`string "%s" is longer than max length (%d)`,
					shortenString(str, defaultMaxPrintLength),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *StringValidator) HasLength(l int) *StringValidator {
	v.tests = append(v.tests, func(str string) error {
		if len(str) != l {
			return errors.New(
				fmt.Sprintf(
					`string "%s" does not have desired length (%d)`,
					shortenString(str, defaultMaxPrintLength),
					l),
			)
		}

		return nil
	})

	return v
}

func (v *StringValidator) Matches(pattern string) *StringValidator {
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	v.tests = append(v.tests, func(str string) error {
		if !r.MatchString(str) {
			return errors.New(
				fmt.Sprintf(
					`string "%s" does not match expected pattern`,
					shortenString(str, defaultMaxPrintLength),
				),
			)
		}

		return nil
	})

	return v
}
