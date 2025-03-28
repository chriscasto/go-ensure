package ensure

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
