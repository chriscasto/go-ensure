package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"testing"
)

// constructBad generates a function that will call the constructor
// This enables the call to be defined in a test case but called in a context
// that can catch the anticipated panic.
func constructBad[T any](empty T, fields ensure.Fields) func() error {
	return func() error {
		bad := ensure.Struct[T](fields)
		if err := bad.Validate(empty); err != nil {
			return err
		}
		return nil
	}
}

func TestStructValidator_Construct(t *testing.T) {
	// Each of these test cases should result in a panic
	testCases := map[string]struct {
		construct func() error
	}{
		"not struct": {
			construct: constructBad(1, ensure.Fields{
				"foo": ensure.String(),
			}),
		},
		"invalid field": {
			construct: constructBad(testStruct{}, ensure.Fields{
				// Field "foo" does not exist in our struct
				"foo": ensure.String(),
			}),
		},
		"wrong field type": {
			construct: constructBad(testStruct{}, ensure.Fields{
				// This should be int, not string
				"Int": ensure.String(),
			}),
		},
		"wrong number subtype": {
			construct: constructBad(testStruct{}, ensure.Fields{
				// This should be int, not float64
				"Int": ensure.Number[float64](),
			}),
		},
		"wrong number size": {
			construct: constructBad(testStruct{}, ensure.Fields{
				// This should be int, not int8
				"Int": ensure.Number[int8](),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()

			if err := tc.construct(); err != nil {
				t.Errorf("validation occured and generated an error: %s", err.Error())
			}
		})
	}
}

func TestStructValidator_ValidateStruct(t *testing.T) {
	testCases := map[string]struct {
		f         ensure.Fields
		s         testStruct
		expectErr bool
	}{
		"single string expect pass": {
			f: ensure.Fields{
				"Str": ensure.String().HasLength(3),
			},
			s:         testStruct{Str: "foo"},
			expectErr: false,
		},
		"single string expect err": {
			f: ensure.Fields{
				"Str": ensure.String().HasLength(4),
			},
			s:         testStruct{Str: "foo"},
			expectErr: true,
		},
		"single int expect pass": {
			f: ensure.Fields{
				"Int": ensure.Number[int]().IsGreaterThan(1),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: false,
		},
		"single int expect fail": {
			f: ensure.Fields{
				"Int": ensure.Number[int]().IsGreaterThan(10),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: true,
		},
		"single float expect pass": {
			f: ensure.Fields{
				"Float": ensure.Number[float64]().InRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.0,
			},
			expectErr: false,
		},
		"single float expect err": {
			f: ensure.Fields{
				"Float": ensure.Number[float64]().InRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.2,
			},
			expectErr: true,
		},
		"multiple fields expect pass": {
			f: ensure.Fields{
				"Str": ensure.String().HasLength(3),
				"Int": ensure.Number[int]().IsGreaterThan(0),
			},
			s: testStruct{
				Str: "foo",
				Int: 3,
			},
			expectErr: false,
		},
		"multiple fields expect err": {
			f: ensure.Fields{
				"Str": ensure.String().HasLength(3),
				"Int": ensure.Number[int]().IsGreaterThan(0),
			},
			s: testStruct{
				Str: "quux",
				Int: 0,
			},
			expectErr: true,
		},
		"all fields expect pass": {
			f: ensure.Fields{
				"Str":   ensure.String().HasLength(3),
				"Int":   ensure.Number[int]().IsGreaterThan(0),
				"Float": ensure.Number[float64]().IsLessThan(4.2),
			},
			s: testStruct{
				Str:   "foo",
				Int:   3,
				Float: 4.1,
			},
			expectErr: false,
		},
		"all fields expect err": {
			f: ensure.Fields{
				"Str":   ensure.String().HasLength(3),
				"Int":   ensure.Number[int]().IsGreaterThan(0),
				"Float": ensure.Number[float64]().IsLessThan(4.2),
			},
			s: testStruct{
				Str:   "quux",
				Int:   0,
				Float: 4.3,
			},
			expectErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v := ensure.Struct[testStruct](tc.f)
			err := v.ValidateStruct(tc.s)
			if err != nil && !tc.expectErr {
				t.Errorf("Struct().Validate(); expected no error, got %s", err)
			} else if err == nil && tc.expectErr {
				t.Errorf("Struct().Validate(); expected error but got none")
			}
		})
	}
}

func TestStructValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Struct[testStruct](ensure.Fields{}))
}
