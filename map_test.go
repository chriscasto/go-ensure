package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"testing"
)

type mapTestCase[K comparable, V any] struct {
	vals     map[K]V
	willPass bool
}

type mapTestCases[K comparable, V any] map[string]mapTestCase[K, V]

func (tcs mapTestCases[K, V]) run(t *testing.T, mv *ensure.MapValidator[K, V], method string) {
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := mv.Validate(tc.vals)
			if err != nil && tc.willPass {
				t.Errorf(`Map().%s.Validate(%v); expected no error, got "%s"`, method, tc.vals, err)
			} else if err == nil && !tc.willPass {
				t.Errorf(`Map().%s.Validate(%v); expected error but got none`, method, tc.vals)
			}
		})
	}
}

func testMapType[K comparable, V any](t *testing.T, name string, expect string) {
	t.Run(name, func(t *testing.T) {
		av := ensure.Map[K, V]()
		vType := av.Type()
		if vType != expect {
			t.Errorf("Map.Type() = %s; want %s", vType, expect)
		}
	})
}

func TestMapValidator_Type(t *testing.T) {
	testMapType[string, string](t, "string => string", "map[string]string")
	testMapType[string, int](t, "string => int", "map[string]int")
	testMapType[string, float64](t, "string => float64", "map[string]float64")
	testMapType[int, string](t, "int => string", "map[int]string")
	testMapType[int, int](t, "int => int", "map[int]int")
	testMapType[string, testStruct](t, "string => testStruct", "map[string]ensure_test.testStruct")
	testMapType[string, []int](t, "string => []int", "map[string][]int")
}

func TestMapValidator_IsNotEmpty(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, false},
		"one":   {map[string]int{"one": 1}, true},
	}

	testCases.run(t, ensure.Map[string, int]().IsNotEmpty(), "IsNotEmpty()")
}

func TestMapValidator_Length_Equals(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, false},
		"one":   {map[string]int{"one": 1}, true},
		"two":   {map[string]int{"one": 1, "two": 2}, false},
	}

	count := 1
	testCases.run(
		t,
		ensure.Map[string, int]().HasLengthWhere(ensure.Length().Equals(count)),
		fmt.Sprintf("HasLengthWhere.Equals(%d)", count),
	)
}

func TestMapValidator_HasCount(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, false},
		"one":   {map[string]int{"one": 1}, true},
		"two":   {map[string]int{"one": 1, "two": 2}, false},
	}

	count := 1
	testCases.run(
		t,
		ensure.Map[string, int]().HasCount(count),
		fmt.Sprintf("HasCount(%d)", count),
	)
}

func TestMapValidator_HasFewerThan(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, true},
		"one":   {map[string]int{"one": 1}, true},
		"two":   {map[string]int{"one": 1, "two": 2}, false},
	}

	count := 2
	testCases.run(
		t,
		ensure.Map[string, int]().HasFewerThan(count),
		fmt.Sprintf("HasFewerThan(%d)", count),
	)
}

func TestMapValidator_HasMoreThan(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, false},
		"one":   {map[string]int{"one": 1}, false},
		"two":   {map[string]int{"one": 1, "two": 2}, true},
	}

	count := 1
	testCases.run(
		t,
		ensure.Map[string, int]().HasMoreThan(count),
		fmt.Sprintf("HasMoreThan(%d)", count),
	)
}

func TestMapValidator_EachKey(t *testing.T) {
	t.Run("panic if validator type doesn't match key type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Map[string, int]().EachKey(ensure.Number[int]())

		if err := bad.Validate(""); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}

	})

	testCases := mapTestCases[string, int]{
		"one good": {
			vals: map[string]int{
				"abcd": 4,
			},
			willPass: true,
		},
		"all good": {
			vals: map[string]int{
				"abcd": 4,
				"wxyz": 4,
			},
			willPass: true,
		},
		"one bad": {
			vals: map[string]int{
				"abcd": 4,
				"a":    1,
			},
			willPass: false,
		},
	}

	validMap := ensure.Map[string, int]().EachKey(
		ensure.String().IsLongerThan(3),
	)

	testCases.run(t, validMap, "EachKey()")
}

func TestMapValidator_EachValue(t *testing.T) {
	t.Run("panic if validator type doesn't match value type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		bad := ensure.Map[string, int]().EachValue(ensure.String())

		if err := bad.Validate(""); err != nil {
			t.Errorf("validation occured and generated an error: %s", err.Error())
		}
	})

	testCases := mapTestCases[string, int]{
		"one good": {
			vals: map[string]int{
				"abcd": 4,
			},
			willPass: true,
		},
		"all good": {
			vals: map[string]int{
				"abcd": 4,
				"wxyz": 4,
			},
			willPass: true,
		},
		"one bad": {
			vals: map[string]int{
				"abcd": 4,
				"a":    1,
			},
			willPass: false,
		},
	}

	validMap := ensure.Map[string, int]().EachValue(
		ensure.Number[int]().IsGreaterThan(2),
	)

	testCases.run(t, validMap, "EachValue()")
}

func TestMapValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Map[string, int]())
	runDefaultValidatorTestCases(t, ensure.Map[string, []int]())
	runDefaultValidatorTestCases(t, ensure.Map[string, testStruct]())
}
