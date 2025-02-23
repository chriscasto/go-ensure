package main

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
)

func main() {
	// map should not be empty but have fewer than 5 values
	m := ensure.Map[string, int]().IsNotEmpty().HasFewerThan(5).EachValue(
		// each value should be an integer < 10
		ensure.Number[int]().IsLessThan(10),
	)

	good := map[string]int{"m": 1, "b": 2, "c": 3, "d": 4}

	if err := m.Validate(good); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooLong := map[string]int{"m": 1, "b": 2, "c": 3, "d": 4, "e": 5}

	if err := m.Validate(tooLong); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooLarge := map[string]int{"m": 11, "b": 22, "c": 33, "d": 44}

	if err := m.Validate(tooLarge); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}
}
