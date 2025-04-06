package with

// UntypedValidator can accept an arbitrary value for validation
type UntypedValidator interface {
	// Type returns the type of value this UntypedValidator is able to validate
	Type() string

	// ValidateUntyped attempts to validate a value of unknown type
	ValidateUntyped(any) error
}

// Validator can validate against a specified type
type Validator[T any] interface {
	// UntypedValidator methods are included
	UntypedValidator

	// Validate runs any checks against the passed value
	Validate(T) error
}

// Validators is a helper type for defining field and method validators for structs
type Validators map[string]UntypedValidator

// DisplayNames provides a mapping from struct field name to human-understandable name
// Example: "FirstName" => "First Name", "Dob" => "Date of Birth"
type DisplayNames map[string]string
