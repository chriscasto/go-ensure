package ensure

import (
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"iter"
	"reflect"
)

// Length is a convenience function for creating a validator to be used on types with a length property
func Length() *NumberValidator[int] {
	return Number[int]()
}

// newTypeErrorFromTypes is a helper function for creating a TypeError in the
// common scenario where you have the names of the type you want and the type
// you have already available in string form
func newTypeErrorFromTypes(want string, got string) *TypeError {
	return NewTypeError(fmt.Sprintf(
		`expected "%s"; got "%s"`, want, got,
	))
}

// testType compares a value against an expected type and returns a type error if they don't match
func testType(value any, expect string) *TypeError {
	valType := reflect.TypeOf(value).String()

	if valType != expect {
		return newTypeErrorFromTypes(expect, valType)
	}

	return nil
}

func getValidationOptions(options []*with.ValidationOptions) *with.ValidationOptions {
	if len(options) > 0 {
		return options[0]
	}

	// default options
	return with.Options()
}

// checkFunc defines a function that performs a check against a value
type checkFunc[T any] func(T) error

// valChecks is a collection of checkFunc functions
type valChecks[T any] []checkFunc[T]

// newValidationChecks creates a new valChecks struct with an optional set of initial checkFunc functions
func newValidationChecks[T any](initChecks ...checkFunc[T]) *valChecks[T] {
	vc := make(valChecks[T], len(initChecks))

	for i, chk := range initChecks {
		vc[i] = chk
	}

	return &vc
}

// Append adds another checkFunc to a valChecks collection
func (vc *valChecks[T]) Append(check checkFunc[T]) {
	*vc = append(*vc, check)
}

// Evaluate runs every checkFunc against a value and returns the first error
// Additional checkFuncs can be added at call time
func (vc *valChecks[T]) Evaluate(target T, addlChecks ...checkFunc[T]) error {
	// handle stored checks first
	for _, fn := range *vc {
		if err := fn(target); err != nil {
			return err
		}
	}

	// then evaluate any additional checks passed at call time, if any
	for _, fn := range addlChecks {
		// ignore empty values
		if fn == nil {
			continue
		}

		if err := fn(target); err != nil {
			return err
		}
	}

	return nil
}

// EvaluateAll runs every checkFunc against a value and collects all errors returned in a ValidationErrors struct
// Additional checkFuncs can be added at call time
func (vc *valChecks[T]) EvaluateAll(target T, addlChecks ...func(T) error) *ValidationErrors {
	vErrs := newValidationErrors()

	// handle stored checks first
	for _, fn := range *vc {
		if err := fn(target); err != nil {
			vErrs.Append(err)
		}
	}

	// then evaluate any additional checks passed at call time, if any
	for _, fn := range addlChecks {
		// ignore empty values
		if fn == nil {
			continue
		}

		if err := fn(target); err != nil {
			vErrs.Append(err)
		}
	}

	// always return vErrs, even if none are added
	// nil checks against interfaces in Go are stupidly error-prone, so need this to avoid segfaults
	return vErrs
}

// seqChecks is a collection of checkFunc functions to be valuated against a sequence
type seqChecks[T any] struct {
	c *valChecks[T]
}

// newSeqChecks creates a new seqChecks collection with an optional initial set of checkFunc functions
func newSeqChecks[T any](initChecks ...checkFunc[T]) *seqChecks[T] {
	return &seqChecks[T]{
		c: newValidationChecks[T](initChecks...),
	}
}

// Append adds another checkFunc to a seqChecks collection
func (sc *seqChecks[T]) Append(check func(T) error) {
	sc.c.Append(check)
}

// Evaluate runs every checkFunc against a value and returns the first error
func (sc *seqChecks[T]) Evaluate(seq iter.Seq[T]) error {
	for v := range seq {
		if err := sc.c.Evaluate(v); err != nil {
			return err
		}
	}

	return nil
}

// EvaluateAll runs every checkFunc against each value in a sequence and collects
// all errors returned in a ValidationErrors struct
func (sc *seqChecks[T]) EvaluateAll(seq iter.Seq[T]) *ValidationErrors {
	vErrs := newValidationErrors()

	for v := range seq {
		if err := sc.c.EvaluateAll(v); err != nil {
			vErrs.Append(err)
		}
	}

	return vErrs
}
