package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

func ptrTo(v any) *any {
	return &v
}

func TestPointer_Type(t *testing.T) {
	testCases := map[string]struct {
		v            with.Validator
		expectedType string
	}{
		"string": {ensure.String(), "*string"},
		"int":    {ensure.Number[int](), "*int"},
		"float":  {ensure.Number[float64](), "*float64"},
		"map":    {ensure.Map[string, bool](), "*map[string]bool"},
		"struct": {ensure.Struct[testStruct](), "*ensure_test.testStruct"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ptr := ensure.Pointer(tc.v)

			if ptr.Type() != tc.expectedType {
				t.Errorf(`Pointer().Type(); expected %s, got "%s"`, tc.expectedType, ptr.Type())
			}
		})
	}
}

func TestPointer_Validate(t *testing.T) {
	testCases := map[string]struct {
		v        with.Validator
		value    any
		willPass bool
	}{
		"nil":                {ensure.String(), nil, false},
		"valid string ptr":   {ensure.String().Contains("string"), ptrTo("string"), true},
		"invalid string ptr": {ensure.String().Contains("string"), ptrTo("str"), false},
		"string val":         {ensure.String().Contains("string"), "string", false},
		"valid int ptr":      {ensure.Number[int]().IsOdd(), ptrTo(123), true},
		"invalid int ptr":    {ensure.Number[int]().IsOdd(), ptrTo(12), false},
		"int val":            {ensure.Number[int]().IsOdd(), 123, false},
		"valid float ptr":    {ensure.Number[float64]().IsPositive(), ptrTo(1.0), true},
		"invalid float ptr":  {ensure.Number[float64]().IsPositive(), ptrTo(-1.0), false},
		"float val":          {ensure.Number[float64]().IsPositive(), 1.0, false},
		"valid map ptr":      {ensure.Map[string, bool]().HasCount(1), ptrTo(map[string]bool{"abc": true}), true},
		"invalid map ptr":    {ensure.Map[string, bool]().HasCount(1), ptrTo(map[string]bool{}), false},
		"map val":            {ensure.Map[string, bool]().HasCount(1), map[string]bool{"abc": true}, false},
		"struct ptr":         {ensure.Struct[testStruct](), ptrTo(testStruct{}), true},
		"struct val":         {ensure.Struct[testStruct](), testStruct{}, false},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ptr := ensure.Pointer(tc.v)

			err := ptr.Validate(tc.value)
			if err != nil && tc.willPass {
				t.Errorf(`Pointer().Validate(%v); expected no error, got "%s"`, tc.value, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Pointer().Validate(%v); expected error but got none`, tc.value)
			}
		})
	}
}

func TestOptionalPointer(t *testing.T) {
	var nilStr *string
	validStr := ensure.String()

	// This should cause an error when dereferencing the nil pointer
	t.Run("required pointer", func(t *testing.T) {
		ptr := ensure.Pointer(validStr)

		if err := ptr.Validate(nilStr); err == nil {
			t.Errorf(`expected error but got none`)
		}
	})

	// This should not cause an error because validation is optional
	t.Run("optional pointer", func(t *testing.T) {
		ptr := ensure.OptionalPointer(validStr)

		if err := ptr.Validate(nilStr); err != nil {
			t.Errorf(`expected no error but got "%s"`, err.Error())
		}
	})

	// Now try it with a struct field
	type testStructWithNil struct {
		Str *string
	}

	nilStruct := testStructWithNil{}

	// This should cause an error when dereferencing the nil pointer
	t.Run("required pointer", func(t *testing.T) {
		reqPtr := ensure.Struct[testStructWithNil]().HasFields(with.Validators{
			"Str": ensure.Pointer(ensure.String()),
		})

		if err := reqPtr.Validate(nilStruct); err == nil {
			t.Errorf(`expected error but got none`)
		}
	})

	// This should not cause an error because validation is optional
	t.Run("optional pointer", func(t *testing.T) {
		optPtr := ensure.Struct[testStructWithNil]().HasFields(with.Validators{
			"Str": ensure.OptionalPointer(ensure.String()),
		})

		if err := optPtr.Validate(nilStruct); err != nil {
			t.Errorf(`expected no error but got "%s"`, err.Error())
		}
	})
}

func TestPointer_ArrayOfPointers(t *testing.T) {
	one := "one"
	abc := "abc"

	strPtrs := []*string{&one, &abc}

	validArr := ensure.Array[*string]().Each(
		ensure.Pointer(
			ensure.String(),
		),
	)

	// This should not cause an error
	t.Run("array of pointers", func(t *testing.T) {
		if err := validArr.Validate(strPtrs); err != nil {
			t.Errorf(`expected no error but got "%s"`, err.Error())
		}
	})
}

func TestPointer_Nested(t *testing.T) {
	str := "foo"
	pStr := &str
	ppStr := &pStr

	// This should not cause an error
	t.Run("pointer of pointer", func(t *testing.T) {
		// Pointer of a pointer
		ptr := ensure.Pointer(
			ensure.Pointer(ensure.String()),
		)

		if err := ptr.Validate(ppStr); err != nil {
			t.Errorf(`expected no error but got "%s"`, err.Error())
		}
	})
}
