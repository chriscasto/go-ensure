package main

import (
	"fmt"
	ensure "github.com/chriscasto/go-ensure"
)

func main() {
	// number should be an integer < 10
	i := ensure.Number[int]().IsLessThan(10)

	goodInt := 5

	if err := i.Validate(goodInt); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooBigInt := 100

	if err := i.Validate(tooBigInt); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	// number should be a float between 1.0 and 10.0
	f := ensure.Number[float64]().InRange(1.0, 2.0)

	goodFloat := 1.2345

	if err := f.Validate(goodFloat); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooBigFLoat := 2.3456

	if err := f.Validate(tooBigFLoat); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooSmallFloat := 0.1234

	if err := f.Validate(tooSmallFloat); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}
}
