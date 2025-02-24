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

type isEvenTestCases[T ensure.NumberType] struct {
	validator *ensure.NumberValidator[T]
	tests     numTestCases[T]
}

func (tcs isEvenTestCases[T]) run(t *testing.T, method string) {
	tcs.tests.run(t, tcs.validator, method)
}

func makeIsEvenTestCases[T ensure.NumberType](expectEven bool) isEvenTestCases[T] {
	validator := ensure.Number[T]()

	if expectEven {
		validator.IsEven()
	} else {
		validator.IsOdd()
	}

	return isEvenTestCases[T]{
		validator,
		numTestCases[T]{
			"zero": {0, expectEven},
			"one":  {1, !expectEven},
			"two":  {2, expectEven},
		},
	}
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

func TestNumberValidator_Equals(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, false},
		"equal to":     {target, true},
		"greater than": {target + 1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().Equals(target),
		fmt.Sprintf("Equals(%d)", target),
	)
}

func TestNumberValidator_DoesNotEqual(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, true},
		"equal to":     {target, false},
		"greater than": {target + 1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().DoesNotEqual(target),
		fmt.Sprintf("DoesNotEqual(%d)", target),
	)
}

func TestNumberValidator_IsLessThan(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, true},
		"equal to":     {target, false},
		"greater than": {target + 1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsLessThan(target),
		fmt.Sprintf("IsLessThan(%d)", target),
	)
}

func TestNumberValidator_IsLessThanOrEqualTo(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, true},
		"equal to":     {target, true},
		"greater than": {target + 1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsLessThanOrEqualTo(target),
		fmt.Sprintf("IsLessThanOrEqualTo(%d)", target),
	)
}

func TestNumberValidator_IsGreaterThan(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, false},
		"equal to":     {target, false},
		"greater than": {target + 1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsGreaterThan(target),
		fmt.Sprintf("IsGreaterThan(%d)", target),
	)
}

func TestNumberValidator_IsGreaterThanOrEqualTo(t *testing.T) {
	target := 10

	testCases := numTestCases[int]{
		"less than":    {target - 1, false},
		"equal to":     {target, true},
		"greater than": {target + 1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsGreaterThanOrEqualTo(target),
		fmt.Sprintf("IsGreaterThanOrEqualTo(%d)", target),
	)
}

func Test_IsEven(t *testing.T) {
	// most cases are covered in the Number().IsEven()/IsOdd() tests below, so
	// just covering edge cases here

	// this should be false because 2.1 is not a whole number
	if ensure.IsEven[float32]("float32", 2.1) {
		t.Errorf("expect isEven() to return false for %v", 2.1)
	}

	// this should be false because 2.1 is not a whole number
	if ensure.IsEven[float64]("float64", 2.1) {
		t.Errorf("expect isEven() to return false for %v", 2.1)
	}

	t.Run("panic if type is not a number", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// this should panic because the type string passed is not a valid number
		ensure.IsEven[int]("string", 2)
	})
}

func TestNumberValidator_IsEven(t *testing.T) {
	method := "IsEven()"
	expectEven := true
	makeIsEvenTestCases[int](expectEven).run(t, method)
	makeIsEvenTestCases[int8](expectEven).run(t, method)
	makeIsEvenTestCases[int16](expectEven).run(t, method)
	makeIsEvenTestCases[int32](expectEven).run(t, method)
	makeIsEvenTestCases[int64](expectEven).run(t, method)
	makeIsEvenTestCases[uint](expectEven).run(t, method)
	makeIsEvenTestCases[uint8](expectEven).run(t, method)
	makeIsEvenTestCases[uint16](expectEven).run(t, method)
	makeIsEvenTestCases[uint32](expectEven).run(t, method)
	makeIsEvenTestCases[uint64](expectEven).run(t, method)
	makeIsEvenTestCases[float32](expectEven).run(t, method)
	makeIsEvenTestCases[float64](expectEven).run(t, method)
}

func TestNumberValidator_IsOdd(t *testing.T) {
	method := "IsOdd()"
	expectEven := false
	makeIsEvenTestCases[int](expectEven).run(t, method)
	makeIsEvenTestCases[int8](expectEven).run(t, method)
	makeIsEvenTestCases[int16](expectEven).run(t, method)
	makeIsEvenTestCases[int32](expectEven).run(t, method)
	makeIsEvenTestCases[int64](expectEven).run(t, method)
	makeIsEvenTestCases[uint](expectEven).run(t, method)
	makeIsEvenTestCases[uint8](expectEven).run(t, method)
	makeIsEvenTestCases[uint16](expectEven).run(t, method)
	makeIsEvenTestCases[uint32](expectEven).run(t, method)
	makeIsEvenTestCases[uint64](expectEven).run(t, method)
	makeIsEvenTestCases[float32](expectEven).run(t, method)
	makeIsEvenTestCases[float64](expectEven).run(t, method)
}

func TestNumberValidator_IsPositive(t *testing.T) {
	testCases := numTestCases[int]{
		"negative one": {-1, false},
		"zero":         {0, false},
		"one":          {1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsPositive(),
		"IsPositive()",
	)
}

func TestNumberValidator_IsNegative(t *testing.T) {
	testCases := numTestCases[int]{
		"negative one": {-1, true},
		"zero":         {0, false},
		"one":          {1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsNegative(),
		"IsNegative()",
	)
}

func TestNumberValidator_IsZero(t *testing.T) {
	testCases := numTestCases[int]{
		"negative one": {-1, false},
		"zero":         {0, true},
		"one":          {1, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsZero(),
		"IsZero()",
	)
}

func TestNumberValidator_IsNotZero(t *testing.T) {
	testCases := numTestCases[int]{
		"negative one": {-1, true},
		"zero":         {0, false},
		"one":          {1, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsNotZero(),
		"IsNotZero()",
	)
}

func TestNumberValidator_IsOneOf(t *testing.T) {
	arr := []int{1, 3, 5}

	testCases := numTestCases[int]{
		"one":   {1, true},
		"two":   {2, false},
		"three": {3, true},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsOneOf(arr),
		fmt.Sprintf("IsOneOf(%v)", arr),
	)
}

func TestNumberValidator_IsNotOneOf(t *testing.T) {
	arr := []int{1, 3, 5}

	testCases := numTestCases[int]{
		"one":   {1, false},
		"two":   {2, true},
		"three": {3, false},
	}

	testCases.run(
		t,
		ensure.Number[int]().IsNotOneOf(arr),
		fmt.Sprintf("IsNotOneOf(%v)", arr),
	)
}

func TestNumberValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Number[int]())
	runDefaultValidatorTestCases(t, ensure.Number[float64]())
}
