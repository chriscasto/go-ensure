package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"testing"
)

type arrayTestCase[T any] struct {
	vals     []T
	willPass bool
}

type arrayDummyStruct struct {
	Foo string
	Bar int
}

func testArrayType[T any](t *testing.T, name string, expect string) {
	t.Run(name, func(t *testing.T) {
		av := ensure.Array[T]()
		vType := av.Type()
		if vType != expect {
			t.Errorf("Array.Type() = %s; want %s", vType, expect)
		}
	})
}

type arrayTestCases[T any] map[string]arrayTestCase[T]

func (tcs arrayTestCases[T]) run(t *testing.T, av *ensure.ArrayValidator[T], method string) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := av.ValidateArray(tc.vals)
			if err != nil && tc.willPass {
				t.Errorf(`Array().%s.Validate(%v); expected no error, got "%s"`, method, tc.vals, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Array().%s.Validate(%v); expected error but got none`, method, tc.vals)
			}
		})
	}
}

func TestArrayValidator_Type(t *testing.T) {
	testArrayType[string](t, "string", "[]string")
	testArrayType[int](t, "int", "[]int")
	testArrayType[float64](t, "float", "[]float64")
	testArrayType[arrayDummyStruct](t, "struct", "[]ensure_test.arrayDummyStruct")
	testArrayType[[]int](t, "array of int", "[][]int")
}

func TestArrayValidator_IsEmpty(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, true},
		"one":   {[]int{1}, false},
	}

	testCases.run(t, ensure.Array[int]().IsEmpty(), "IsEmpty()")
}

func TestArrayValidator_IsNotEmpty(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, false},
		"one":   {[]int{1}, true},
	}

	testCases.run(t, ensure.Array[int]().IsNotEmpty(), "IsNotEmpty()")
}

func TestArrayValidator_Length_Equals(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, false},
		"one":   {[]int{1}, true},
		"two":   {[]int{1, 2}, false},
	}

	count := 1
	testCases.run(
		t,
		ensure.Array[int]().HasLengthWhere(ensure.Length().Equals(count)),
		fmt.Sprintf("HasLengthWhere(Length().Equals(%d))", count),
	)
}

func TestArrayValidator_HasCount(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, false},
		"one":   {[]int{1}, true},
		"two":   {[]int{1, 2}, false},
	}

	count := 1
	testCases.run(
		t,
		ensure.Array[int]().HasCount(count),
		fmt.Sprintf("HasCount(%d)", count),
	)
}

func TestArrayValidator_HasFewerThan(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, true},
		"one":   {[]int{1}, false},
		"two":   {[]int{1, 2}, false},
	}

	count := 1
	testCases.run(
		t,
		ensure.Array[int]().HasFewerThan(count),
		fmt.Sprintf("HasFewerThan(%d)", count),
	)
}

func TestArrayValidator_HasMoreThan(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty": {[]int{}, false},
		"one":   {[]int{1}, false},
		"two":   {[]int{1, 2}, true},
	}

	count := 1
	testCases.run(
		t,
		ensure.Array[int]().HasMoreThan(count),
		fmt.Sprintf("HasMoreThan(%d)", count),
	)
}

func TestArrayValidator_Each(t *testing.T) {
	testCases := arrayTestCases[int]{
		"empty":      {[]int{}, true}, // passes due to no individual failures
		"none pass":  {[]int{1}, false},
		"one passes": {[]int{1, 2}, false},
		"all pass":   {[]int{2, 3, 4}, true},
	}

	valMin := 1
	testCases.run(
		t,
		ensure.Array[int]().Each(ensure.Number[int]().IsGreaterThan(valMin)),
		fmt.Sprintf("Each(IsGreaterThan(%d))", valMin),
	)
}

func TestArrayValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Array[bool]())
	runDefaultValidatorTestCases(t, ensure.Array[string]())
	runDefaultValidatorTestCases(t, ensure.Array[int]())
	runDefaultValidatorTestCases(t, ensure.Array[float64]())
	runDefaultValidatorTestCases(t, ensure.Array[testStruct]())
}
