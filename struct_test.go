package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
	"time"
)

// constructBad generates a function that will call the constructor
// This enables the call to be defined in a test case but called in a context
// that can catch the anticipated panic.
func constructBad[T any](empty T) func() error {
	return func() error {
		bad := ensure.Struct[T]()
		if err := bad.Validate(empty); err != nil {
			return err
		}
		return nil
	}
}

func setBadFields[T any](empty T, fields with.Validators) func() error {
	return func() error {
		bad := ensure.Struct[T]().HasFields(fields)
		if err := bad.Validate(empty); err != nil {
			return err
		}
		return nil
	}
}

func setBadGetters[T any](empty T, methods with.Validators) func() error {
	return func() error {
		bad := ensure.Struct[T]().HasGetters(methods)
		if err := bad.Validate(empty); err != nil {
			return err
		}
		return nil
	}
}

type structTestCase[T any] struct {
	val      T
	willPass bool
}

type structTestCases[T any] map[string]structTestCase[T]

func (tcs structTestCases[T]) run(t *testing.T, av *ensure.StructValidator[T], method string) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := av.Validate(tc.val)
			if err != nil && tc.willPass {
				t.Errorf(`Struct().%s.Validate(%v); expected no error, got "%s"`, method, tc.val, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Struct().%s.Validate(%v); expected error but got none`, method, tc.val)
			}
		})
	}
}

func TestStructValidator_Construct(t *testing.T) {
	t.Run("not struct", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		if err := constructBad(1)(); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}
	})
}

func TestStructValidator_HasFields(t *testing.T) {
	// Each of these test cases should result in a panic
	panicTestCases := map[string]struct {
		construct func() error
	}{
		"invalid field": {
			construct: setBadFields(testStruct{}, with.Validators{
				// Field "foo" does not exist in our struct
				"foo": ensure.String(),
			}),
		},
		"wrong field type": {
			construct: setBadFields(testStruct{}, with.Validators{
				// This should be int, not string
				"Int": ensure.String(),
			}),
		},
		"wrong number subtype": {
			construct: setBadFields(testStruct{}, with.Validators{
				// This should be int, not float64
				"Int": ensure.Number[float64](),
			}),
		},
		"wrong number size": {
			construct: setBadFields(testStruct{}, with.Validators{
				// This should be int, not int8
				"Int": ensure.Number[int8](),
			}),
		},
	}

	for name, tc := range panicTestCases {
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

	evalTestCases := map[string]struct {
		f         with.Validators
		s         testStruct
		expectErr bool
	}{
		"single string expect pass": {
			f: with.Validators{
				"Str": ensure.String().HasLength(3),
			},
			s:         testStruct{Str: "foo"},
			expectErr: false,
		},
		"single string expect err": {
			f: with.Validators{
				"Str": ensure.String().HasLength(4),
			},
			s:         testStruct{Str: "foo"},
			expectErr: true,
		},
		"single int expect pass": {
			f: with.Validators{
				"Int": ensure.Number[int]().IsGreaterThan(1),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: false,
		},
		"single int expect fail": {
			f: with.Validators{
				"Int": ensure.Number[int]().IsGreaterThan(10),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: true,
		},
		"single float expect pass": {
			f: with.Validators{
				"Float": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.0,
			},
			expectErr: false,
		},
		"single float expect err": {
			f: with.Validators{
				"Float": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.2,
			},
			expectErr: true,
		},
		"multiple fields expect pass": {
			f: with.Validators{
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
			f: with.Validators{
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
			f: with.Validators{
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
			f: with.Validators{
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

	for name, tc := range evalTestCases {
		t.Run(name, func(t *testing.T) {
			v := ensure.Struct[testStruct]().HasFields(tc.f)
			err := v.Validate(tc.s)
			if err != nil && !tc.expectErr {
				t.Errorf("Struct().Validate(); expected no error, got %s", err)
			} else if err == nil && tc.expectErr {
				t.Errorf("Struct().Validate(); expected error but got none")
			}
		})
	}
}

func TestStructValidator_HasGetters(t *testing.T) {
	// Each of these test cases should result in a panic
	panicTestCases := map[string]struct {
		construct func() error
	}{
		"invalid getter": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// Method "Nonexistent" does not exist in our struct
				"Nonexistent": ensure.String(),
			}),
		},
		"wrong return type": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// This should be int, not string
				"GetInt": ensure.String(),
			}),
		},
		"wrong number subtype": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// This should be int, not float64
				"GetInt": ensure.Number[float64](),
			}),
		},
		"wrong number size": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// This should be int, not int8
				"GetInt": ensure.Number[int8](),
			}),
		},
		"too many args": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// This method returns both string and error
				"GetStrWithArg": ensure.String(),
			}),
		},
		"too many return values": {
			construct: setBadGetters(testStruct{}, with.Validators{
				// This method returns both string and error
				"GetStrWithError": ensure.String(),
			}),
		},
	}

	for name, tc := range panicTestCases {
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

	evalTestCases := map[string]struct {
		validators with.Validators
		s          testStruct
		expectErr  bool
	}{
		"single string expect pass": {
			validators: with.Validators{
				"GetStr": ensure.String().HasLength(3),
			},
			s:         testStruct{Str: "foo"},
			expectErr: false,
		},
		"single string expect err": {
			validators: with.Validators{
				"GetStr": ensure.String().HasLength(4),
			},
			s:         testStruct{Str: "foo"},
			expectErr: true,
		},
		"single int expect pass": {
			validators: with.Validators{
				"GetInt": ensure.Number[int]().IsGreaterThan(1),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: false,
		},
		"single int expect fail": {
			validators: with.Validators{
				"GetInt": ensure.Number[int]().IsGreaterThan(10),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: true,
		},
		"single float expect pass": {
			validators: with.Validators{
				"GetFloat": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.0,
			},
			expectErr: false,
		},
		"single float expect err": {
			validators: with.Validators{
				"GetFloat": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.2,
			},
			expectErr: true,
		},
		"multiple fields expect pass": {
			validators: with.Validators{
				"GetStr": ensure.String().HasLength(3),
				"GetInt": ensure.Number[int]().IsGreaterThan(0),
			},
			s: testStruct{
				Str: "foo",
				Int: 3,
			},
			expectErr: false,
		},
		"multiple fields expect err": {
			validators: with.Validators{
				"GetStr": ensure.String().HasLength(3),
				"GetInt": ensure.Number[int]().IsGreaterThan(0),
			},
			s: testStruct{
				Str: "quux",
				Int: 0,
			},
			expectErr: true,
		},
		"all fields expect pass": {
			validators: with.Validators{
				"GetStr":   ensure.String().HasLength(3),
				"GetInt":   ensure.Number[int]().IsGreaterThan(0),
				"GetFloat": ensure.Number[float64]().IsLessThan(4.2),
			},
			s: testStruct{
				Str:   "foo",
				Int:   3,
				Float: 4.1,
			},
			expectErr: false,
		},
		"all fields expect err": {
			validators: with.Validators{
				"GetStr":   ensure.String().HasLength(3),
				"GetInt":   ensure.Number[int]().IsGreaterThan(0),
				"GetFloat": ensure.Number[float64]().IsLessThan(4.2),
			},
			s: testStruct{
				Str:   "quux",
				Int:   0,
				Float: 4.3,
			},
			expectErr: true,
		},
	}

	for name, tc := range evalTestCases {
		t.Run(name, func(t *testing.T) {
			v := ensure.Struct[testStruct]().HasGetters(tc.validators)
			err := v.Validate(tc.s)
			if err != nil && !tc.expectErr {
				t.Errorf("Struct().Validate(); expected no error, got %s", err)
			} else if err == nil && tc.expectErr {
				t.Errorf("Struct().Validate(); expected error but got none")
			}
		})
	}
}

func TestStructValidator_ValidateStruct(t *testing.T) {
	testCases := map[string]struct {
		f         with.Validators
		s         testStruct
		expectErr bool
	}{
		"single string expect pass": {
			f: with.Validators{
				"Str": ensure.String().HasLength(3),
			},
			s:         testStruct{Str: "foo"},
			expectErr: false,
		},
		"single string expect err": {
			f: with.Validators{
				"Str": ensure.String().HasLength(4),
			},
			s:         testStruct{Str: "foo"},
			expectErr: true,
		},
		"single int expect pass": {
			f: with.Validators{
				"Int": ensure.Number[int]().IsGreaterThan(1),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: false,
		},
		"single int expect fail": {
			f: with.Validators{
				"Int": ensure.Number[int]().IsGreaterThan(10),
			},
			s: testStruct{
				Int: 3,
			},
			expectErr: true,
		},
		"single float expect pass": {
			f: with.Validators{
				"Float": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.0,
			},
			expectErr: false,
		},
		"single float expect err": {
			f: with.Validators{
				"Float": ensure.Number[float64]().IsInRange(2.9, 3.1),
			},
			s: testStruct{
				Float: 3.2,
			},
			expectErr: true,
		},
		"multiple fields expect pass": {
			f: with.Validators{
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
			f: with.Validators{
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
			f: with.Validators{
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
			f: with.Validators{
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
			v := ensure.Struct[testStruct]().HasFields(tc.f)
			err := v.Validate(tc.s)
			if err != nil && !tc.expectErr {
				t.Errorf("Struct().Validate(); expected no error, got %s", err)
			} else if err == nil && tc.expectErr {
				t.Errorf("Struct().Validate(); expected error but got none")
			}
		})
	}
}

func TestStructValidator_HasFields_DisplayNames(t *testing.T) {
	t.Run("panic if field name doesn't exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Struct[testStruct]().HasFields(
			with.Validators{
				"Str": ensure.String(),
			},
			with.DisplayNames{
				"String": "String Value",
			},
		)

		if err := bad.Validate(testStruct{}); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}
	})

	validStruct := ensure.Struct[testStruct]().HasFields(
		with.Validators{
			"Str":   ensure.String().HasLength(3),
			"Int":   ensure.Number[int]().IsGreaterThan(0),
			"Float": ensure.Number[float64]().IsLessThan(4.2),
		},
		with.DisplayNames{
			"Str":   "String Value",
			"Int":   "Integer Value",
			"Float": "Decimal Value",
		},
	)

	testCases := map[string]struct {
		val            testStruct
		expectStrInErr string
	}{
		"string err": {
			testStruct{"a", 1, 1.0},
			"String Value",
		},
		"int err": {
			testStruct{"abc", 0, 1.0},
			"Integer Value",
		},
		"float err": {
			testStruct{"abc", 1, 10.0},
			"Decimal Value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validStruct.Validate(tc.val)

			if err == nil {
				t.Errorf("expected error but got none")
			}

			errorChecker := ensure.String().Contains(tc.expectStrInErr)

			if err2 := errorChecker.Validate(err.Error()); err2 != nil {
				t.Errorf(`error should contain alias "%s" but did not (err: "%s")`, tc.expectStrInErr, err)
			}
		})
	}
}

func TestStructValidator_HasGetters_DisplayNames(t *testing.T) {
	t.Run("panic if method doesn't exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Struct[testStruct]().HasGetters(
			with.Validators{
				"GetStr": ensure.String(),
			},
			with.DisplayNames{
				"String": "String Value",
			},
		)

		if err := bad.Validate(testStruct{}); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}
	})

	validStruct := ensure.Struct[testStruct]().HasGetters(
		with.Validators{
			"GetStr":   ensure.String().HasLength(3),
			"GetInt":   ensure.Number[int]().IsGreaterThan(0),
			"GetFloat": ensure.Number[float64]().IsLessThan(4.2),
		},
		with.DisplayNames{
			"GetStr":   "String Value",
			"GetInt":   "Integer Value",
			"GetFloat": "Decimal Value",
		},
	)

	testCases := map[string]struct {
		val            testStruct
		expectStrInErr string
	}{
		"string err": {
			testStruct{"a", 1, 1.0},
			"String Value",
		},
		"int err": {
			testStruct{"abc", 0, 1.0},
			"Integer Value",
		},
		"float err": {
			testStruct{"abc", 1, 10.0},
			"Decimal Value",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validStruct.Validate(tc.val)

			if err == nil {
				t.Errorf("expected error but got none")
			}

			errorChecker := ensure.String().Contains(tc.expectStrInErr)

			if err2 := errorChecker.Validate(err.Error()); err2 != nil {
				t.Errorf(`error should contain alias "%s" but did not (err: "%s")`, tc.expectStrInErr, err)
			}
		})
	}
}

func TestStructValidator_Is(t *testing.T) {
	type Example struct {
		Date time.Time
	}

	notOlderThanSixtyDays := func(date time.Time) error {
		hourMax := 24 * 60

		if time.Since(date).Hours() > float64(hourMax) {
			return fmt.Errorf("time has expired")
		}

		return nil
	}

	testCases := structTestCases[Example]{
		"now": {
			Example{time.Now()},
			true,
		},
		"yesterday": {
			Example{time.Now().AddDate(0, 0, -1)},
			true,
		},
		"90 days ago": {
			Example{time.Now().AddDate(0, 0, -90)},
			false,
		},
	}

	testCases.run(
		t,
		ensure.Struct[Example]().HasFields(with.Validators{
			"Date": ensure.Struct[time.Time]().Is(notOlderThanSixtyDays),
		}),
		"Is()",
	)
}

func TestStructValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Struct[testStruct]())
}
