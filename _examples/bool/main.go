package main

import (
	"fmt"
	"github.com/chriscasto/go-ensure"
	"github.com/chriscasto/go-ensure/with"
)

func main() {
	isTrue := ensure.Bool().IsTrue()
	isFalse := ensure.Bool().IsFalse()

	// Some simple checks demonstrating basic usage

	fmt.Println("Is true true?")
	if err := isTrue.ValidateBool(true); err != nil {
		fmt.Println("  > You should never see this because validation should never fail")
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > true is true!")
	}

	fmt.Println("Is false true?")
	if err := isTrue.ValidateBool(false); err != nil {
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > You should never see this because validation should always fail")
	}

	fmt.Println("Is false false?")
	if err := isFalse.ValidateBool(false); err != nil {
		fmt.Println("  > You should never see this because validation should never fail")
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > false is false!")
	}

	fmt.Println("Is false true?")
	if err := isFalse.ValidateBool(true); err != nil {
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > You should never see this because validation should always fail")
	}

	// A more complex demonstration using Bool() for conditional validation

	type Component struct {
		Enabled bool
		Id      int
		Name    string
	}

	// Validation will fail if the component is enabled
	componentDisabled := ensure.Struct[Component]().HasFields(with.Validators{
		"Enabled": ensure.Bool().IsFalse(),
	})

	// Make sure that Id and Name are not empty
	componentValid := ensure.Struct[Component]().HasFields(with.Validators{
		"Id":   ensure.Number[int]().IsNotZero(),
		"Name": ensure.String().IsNotEmpty(),
	})

	validIfEnabled := ensure.Any(
		componentDisabled, // this will cause validation to succeed and return immediately if Enabled == false
		componentValid,    // this performs the actual validation only if the first check fails
	).WithError("Component is not valid")

	emptyButDisabled := Component{
		Enabled: false,
	}

	emptyButEnabled := Component{
		Enabled: true,
	}

	complete := Component{
		Enabled: true,
		Id:      1,
		Name:    "test",
	}

	fmt.Println("This will validate successfully because the component is disabled")
	if err := validIfEnabled.Validate(emptyButDisabled); err != nil {
		fmt.Println("  > You should never see this because validation should always succeed")
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > Disabled component is valid")
	}

	fmt.Println("This will fail because Id and Name do not have valid values")
	if err := validIfEnabled.Validate(emptyButEnabled); err != nil {
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > You should never see this because validation should always fail")
	}

	fmt.Println("This will validate successfully because all required fields are valid")
	if err := validIfEnabled.Validate(complete); err != nil {
		fmt.Println("  > You should never see this because validation should always succeed")
		fmt.Println(fmt.Sprintf("  > %s", err))
	} else {
		fmt.Println("  > Complete component is valid")
	}
}
