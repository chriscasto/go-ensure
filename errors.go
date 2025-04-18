package ensure

import (
	"errors"
	"fmt"
)

// TypeError indicates a mismatch between the type expected by a validator and
// the type of the value passed to the validator.  This should generally not
// be passed back to the user.
type TypeError struct {
	err string
}

func (e *TypeError) Error() string {
	return e.err
}

func NewTypeError(err string) *TypeError {
	return &TypeError{err: err}
}

// ValidationError is used to indicate a failure while conducting validation
// checks.  These are intended to be safe to return to the user so they can
// correct their input(s)
type ValidationError struct {
	err string
}

func (e *ValidationError) Error() string {
	return e.err
}

// NewValidationError returns a ValidationError with the error message passed to it
func NewValidationError(err string) *ValidationError {
	return &ValidationError{err}
}

// ValidationErrors is a collection of multiple TypeError and ValidationError structs
// It provides a transparent mechanism for returning multiple errors from a validation tree
type ValidationErrors struct {
	tErrs []*TypeError
	vErrs []*ValidationError
}

// newValidationErrors creates and returns a new ValidationErrors struct
func newValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		tErrs: make([]*TypeError, 0),
		vErrs: make([]*ValidationError, 0),
	}
}

// ErrorAsValidationErrors is a helper function for checking if an error is an instance of ValidationErrors
func ErrorAsValidationErrors(err error) *ValidationErrors {
	vErrs := &ValidationErrors{}

	if errors.As(err, &vErrs) {
		return vErrs
	}

	return nil
}

// Append adds an error to the internal list of validation errors
// TypeError is kept separate from ValidationError
// If it's a ValidationErrors error, it gets merged into this one
func (v *ValidationErrors) Append(err error) {

	// The most common case is that it's another set of validation errors from upstream
	ve := &ValidationErrors{}

	if errors.As(err, &ve) {
		v.Extend(ve)
		return
	}

	// If it's a TypeError, keep that separate
	te := &TypeError{}

	if errors.As(err, &te) {
		v.tErrs = append(v.tErrs, te)
		return
	}

	// Otherwise default to adding it as a ValidationError
	v.vErrs = append(v.vErrs, &ValidationError{err: err.Error()})
}

// Extend adds all the errors collected in one ValidationErrors instance into another
func (v *ValidationErrors) Extend(errs *ValidationErrors) {
	// Add type errors from other struct
	for _, err := range errs.tErrs {
		v.tErrs = append(v.tErrs, err)
	}

	// Add validation errors from other struct
	for _, err := range errs.vErrs {
		v.vErrs = append(v.vErrs, err)
	}
}

// HasErrors returns true if there are any errors
func (v *ValidationErrors) HasErrors() bool {
	return v.HasTypeErrors() || v.HasValidationErrors()
}

// HasTypeErrors returns true if there are any TypeErrors
func (v *ValidationErrors) HasTypeErrors() bool {
	return len(v.tErrs) > 0
}

// GetTypeErrors returns all collected TypeErrors
func (v *ValidationErrors) GetTypeErrors() []*TypeError {
	return v.tErrs
}

// HasValidationErrors returns true if there are any ValidationErrors
func (v *ValidationErrors) HasValidationErrors() bool {
	return len(v.vErrs) > 0
}

// GetValidationErrors returns all collected ValidationErrors
func (v *ValidationErrors) GetValidationErrors() []*ValidationError {
	return v.vErrs
}

// Error is the implementation of the "error" interface
func (v *ValidationErrors) Error() string {
	// Return the first validation error by default
	if len(v.vErrs) > 0 {
		return v.vErrs[0].Error()
	}

	// If there are none, return a vague error about available type errors
	// We don't want to return the actual type error since this may expose internal information
	if len(v.tErrs) > 0 {
		// Print a generic message to avoid leaking validation details to the user in case of improper error handling
		// If you want the actual message(s), use HasTypeErrors() and GetTypeErrors()
		return fmt.Sprintf("there were %d type errors", len(v.tErrs))
	}

	// We shouldn't get here, but just in case
	return "there were no validation errors"
}
