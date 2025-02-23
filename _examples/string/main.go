package main

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
)

func main() {
	// string should be less than 7 characters long and only contain alpa characters
	s := ensure.String().IsShorterThan(7).Matches(ensure.Alpha)

	good := "AbC"

	if err := s.Validate(good); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	tooLong := "abcdefg"

	if err := s.Validate(tooLong); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	noMatch := "abc123"

	if err := s.Validate(noMatch); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("OK")
	}
}
