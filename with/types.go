package with

// Validator represents an object that can evaluate a passed value against a set of checks
type Validator interface {
	// Validate runs any checks against the passed value
	Validate(any) error
	// Type returns the type of value this Validator is able to validate
	Type() string
}

// Validators is a helper type for defining field and method validators for structs
type Validators map[string]Validator

// FriendlyNames provides a mapping from struct field name to human-understandable name
// Example: "FirstName" => "First Name", "Dob" => "Date of Birth"
type FriendlyNames map[string]string
