package main

import (
	"fmt"
	ensure "github.com/chriscasto/go-ensure"
)

func main() {
	// array should not be empty but have fewer than 5 values
	a := ensure.Array[int]().IsNotEmpty().HasFewerThan(5).Each(
		// each value should be an integer < 10
		ensure.Number[int]().IsLessThan(10),
	)

	good := []int{1, 2, 3, 4}

	if err := a.Validate(good); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooLong := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if err := a.Validate(tooLong); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooLarge := []int{11, 22, 33, 44}

	if err := a.Validate(tooLarge); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}
}
