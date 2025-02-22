package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"testing"
)

type numTestCase[T ensure.NumberType] struct {
	value    T
	willPass bool
}

type numTestCases[T ensure.NumberType] map[string]numTestCase[T]

func (tcs numTestCases[T]) run(t *testing.T, sv *ensure.NumberValidator[T], method string) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := sv.Validate(tc.value)
			if err != nil && tc.willPass {
				t.Errorf(`Number().%s.Validate(%v); expected no error, got "%s"`, method, tc.value, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Number().%s.Validate(%v); expected error but got none`, method, tc.value)
			}
		})
	}
}

func numTestType[T ensure.NumberType](t *testing.T, name string, expect string) {
	t.Run(name, func(t *testing.T) {
		av := ensure.Number[T]()
		vType := av.Type()
		if vType != expect {
			t.Errorf("Number.Type() = %s; want %s", vType, expect)
		}
	})
}

func TestNumberValidator_Type(t *testing.T) {
	numTestType[int](t, "int", "int")
	numTestType[int8](t, "8-bit int", "int8")
	numTestType[int64](t, "64-bit int", "int64")
	numTestType[uint](t, "unsigned int", "uint")
	numTestType[float64](t, "64-bit float", "float64")
}

func TestNumberValidator_InRange(t *testing.T) {
	rangeMin := 1
	rangeMax := 10

	// Check to make sure the method panics if max < min
	t.Run("panic if max < min", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Number[int]().InRange(rangeMax, rangeMin)
		if err := bad.Validate(rangeMin); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}

	})

	testCases := numTestCases[int]{
		"less than":       {rangeMin - 1, false},
		"bottom of range": {rangeMin, true},
		"top of range":    {rangeMax, false},
		"greater than":    {rangeMax + 1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().InRange(rangeMin, rangeMax),
		fmt.Sprintf("InRange(%d, %d)", rangeMin, rangeMax),
	)
}

func TestNumberValidator_LessThan(t *testing.T) {
	valMax := 10

	testCases := numTestCases[int]{
		"less than":    {valMax - 1, true},
		"equal to":     {valMax, false},
		"greater than": {valMax + 1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsLessThan(valMax),
		fmt.Sprintf("IsLessThan(%d)", valMax),
	)
}

func TestNumberValidator_GreaterThan(t *testing.T) {
	valMin := 10

	testCases := numTestCases[int]{
		"less than":    {valMin - 1, false},
		"equal to":     {valMin, false},
		"greater than": {valMin + 1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsGreaterThan(valMin),
		fmt.Sprintf("IsGreaterThan(%d)", valMin),
	)
}

func TestNumberValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Number[int]())
	runDefaultValidatorTestCases(t, ensure.Number[float64]())
}
