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
// checks.  These are generally safe to return to the user so they can correct
// their input(s)
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

type ValidationErrors struct {
	tErrs []*TypeError
	vErrs []*ValidationError
}

func newValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		tErrs: make([]*TypeError, 0),
		vErrs: make([]*ValidationError, 0),
	}
}

func ErrorAsValidationErrors(err error) *ValidationErrors {
	vErrs := &ValidationErrors{}

	if errors.As(err, &vErrs) {
		return vErrs
	}

	return nil
}

func (v *ValidationErrors) Append(err error) {

	// The most common case is that it's a another set of validation errors from upstream
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

func (v *ValidationErrors) HasErrors() bool {
	return v.HasTypeErrors() || v.HasValidationErrors()
}

func (v *ValidationErrors) HasTypeErrors() bool {
	return len(v.tErrs) > 0
}

func (v *ValidationErrors) GetTypeErrors() []*TypeError {
	return v.tErrs
}

func (v *ValidationErrors) HasValidationErrors() bool {
	return len(v.vErrs) > 0
}

func (v *ValidationErrors) GetValidationErrors() []*ValidationError {
	return v.vErrs
}

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

func test() {

	ve := newValidationErrors()
	out := []error{}

	if ve.HasTypeErrors() {
		for _, err := range ve.GetTypeErrors() {
			fmt.Println(err.Error())
		}

	}

	if ve.HasValidationErrors() {
		for _, err := range ve.GetValidationErrors() {
			out = append(out, err)
		}
	}
}
