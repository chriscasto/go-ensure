package ensure_test

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
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

// TestMapValidator_IsValidator checks to make sure the MapValidator implements the Validator interfaces
func TestMapValidator_IsValidator(t *testing.T) {
	var _ with.UntypedValidator = ensure.Map[string, string]()
	var _ with.Validator[map[string]string] = ensure.Map[string, string]()
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

func TestMapValidator_IsEmpty(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"empty": {map[string]int{}, true},
		"one":   {map[string]int{"one": 1}, false},
	}

	testCases.run(t, ensure.Map[string, int]().IsEmpty(), "IsEmpty()")
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
		fmt.Sprintf("HasLengthWhere(Length().Equals(%d))", count),
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

func TestMapValidator_Is(t *testing.T) {
	testCases := mapTestCases[string, int]{
		"one good": {
			vals: map[string]int{
				"abcd": 4,
			},
			willPass: true,
		},
		"all good": {
			vals: map[string]int{
				"abc":  3,
				"wxyz": 4,
			},
			willPass: true,
		},
		"one bad": {
			vals: map[string]int{
				"abc": 4,
				"a":   1,
			},
			willPass: false,
		},
	}

	equalKeys := func(m map[string]int) error {
		for k, v := range m {
			if len(k) != v {
				return fmt.Errorf("key length (%d) must equal value (%d)", len(k), v)
			}
		}

		return nil
	}

	testCases.run(t, ensure.Map[string, int]().Is(equalKeys), "Is()")
	testCases.run(t, ensure.Map[string, int]().Has(equalKeys), "Has()")
}

func TestMapValidator_MultiError(t *testing.T) {
	type exampleMap = map[string]int

	mapTestCases := multiErrTestCases[exampleMap]{
		"empty": {exampleMap{}, 1},                               // fails not empty
		"one":   {exampleMap{"one": 1}, 0},                       // fails none
		"two":   {exampleMap{"one": 1, "two": 2}, 1},             // fails odd
		"three": {exampleMap{"one": 1, "two": 2, "three": 3}, 3}, // fails fewer than 3, str length 3, odd
		"four":  {exampleMap{"four": 4}, 2},                      // fails str length 3, odd
	}

	mapTestCases.run(t,
		ensure.Map[string, int]().IsNotEmpty().HasLengthWhere(
			ensure.Length().IsLessThan(3),
		).EachKey(
			ensure.String().HasLength(3),
		).EachValue(
			ensure.Number[int]().IsOdd(),
		),
	)
}

func TestMapValidator_Validate(t *testing.T) {
	// see util_test.go
	runDefaultValidatorTestCases(t, ensure.Map[string, int]())
	runDefaultValidatorTestCases(t, ensure.Map[string, []int]())
	runDefaultValidatorTestCases(t, ensure.Map[string, testStruct]())
}
