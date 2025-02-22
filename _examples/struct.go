package main

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
)

type testStruct struct {
	Foo string
	Bar int
	Baz []float64
}

func main() {
	// struct should be type main.testStruct
	s := ensure.Struct[testStruct](ensure.Fields{
		// field Foo should be a string with more than 3 characters
		"Foo": ensure.String().IsLongerThan(3),
		// field Bar should be an integer > 10
		"Bar": ensure.Number[int]().IsGreaterThan(10),
		// field Baz is an array of floats
		"Baz": ensure.Array[float64]().Each(
			// each value should be between 1.0 and 10.0
			ensure.Number[float64]().InRange(1.0, 10.0),
		),
	})

	good := testStruct{
		Foo: "quux",
		Bar: 11,
		Baz: []float64{1.0, 2.0},
	}

	if err := s.Validate(good); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	fooTooShort := testStruct{
		Foo: "a",
		Bar: 11,
		Baz: []float64{1.0, 2.0},
	}

	if err := s.Validate(fooTooShort); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	barTooSmall := testStruct{
		Foo: "quux",
		Bar: 1,
		Baz: []float64{1.0, 2.0},
	}

	if err := s.Validate(barTooSmall); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	bazNotInRange := testStruct{
		Foo: "quux",
		Bar: 11,
		Baz: []float64{0.0, 11.0},
	}

	if err := s.Validate(bazNotInRange); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}
}
