package main

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
)

type testStruct struct {
	Foo string
	Bar int
	baz []float64
}

func (ts *testStruct) GetBaz() []float64 {
	return ts.baz
}

func validateStruct(msg string, v with.Validator, s testStruct) {
	fmt.Printf("%s\n", msg)
	if err := v.Validate(s); err != nil {
		fmt.Printf("->  Error: %v\n", err)
	} else {
		fmt.Println("->  OK")
	}
}

func main() {
	// struct should be type main.testStruct
	s := ensure.Struct[testStruct]()

	// Define validators for fields
	s.HasFields(
		with.Validators{
			// field Foo should be a string with more than 3 characters
			"Foo": ensure.String().IsLongerThan(3),
			// field Bar should be an integer > 10
			"Bar": ensure.Number[int]().IsGreaterThan(10),
		},
		// define some user-friendly aliases for our fields to use when returning errors
		with.DisplayNames{
			"Foo": "FOOOOOO!",
			"Bar": "BAR BAR BAR",
		},
	)

	// Define validators for getter methods
	s.HasGetters(with.Validators{
		// method GetBaz should return an array of floats
		"GetBaz": ensure.Array[float64]().Each(
			// each value should be between 1.0 and 10.0
			ensure.Number[float64]().IsInRange(1.0, 10.0),
		),
	}, with.DisplayNames{
		"GetBaz": "Bazzler",
	})

	validateStruct("This one should pass", s, testStruct{
		Foo: "quux",
		Bar: 11,
		baz: []float64{1.0, 2.0},
	})

	validateStruct("This one should fail because Foo is too short", s, testStruct{
		Foo: "a",
		Bar: 11,
		baz: []float64{1.0, 2.0},
	})

	validateStruct("This one should fail because Bar is too small", s, testStruct{
		Foo: "quux",
		Bar: 1,
		baz: []float64{1.0, 2.0},
	})

	validateStruct(
		"This one should fail because the array returned by GetBaz contains a number outside the expected range",
		s,
		testStruct{
			Foo: "quux",
			Bar: 11,
			baz: []float64{0.0, 11.0},
		})
}
