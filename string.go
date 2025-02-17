package valid

import (
	"errors"
	"fmt"
	"reflect"
)

// define the max length for string contents in an error message
const defaultMaxPrintLength = 30

// shortens a string of arbitrary length to a reasonable value for logging
func shortenString(s string, maxLen int) string {
	if maxLen < 5 {
		maxLen = 5
	}

	if len(s) < maxLen {
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

func (v *StringValidator) Kind() reflect.Kind {
	return reflect.String
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

func (v *StringValidator) OneOf(vals ...string) *StringValidator {
	return v
}

func (v *StringValidator) LongerThan(l int) *StringValidator {

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

func (v *StringValidator) ShorterThan(l int) *StringValidator {
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
