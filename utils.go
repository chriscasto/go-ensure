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

// seqChecks is a collection of checkFunc functions to be valuated against a sequence
type seqChecks[T any] struct {
	c *valChecks[T]
}

// newSeqChecks creates a new seqChecks collection with an optional initial set of checkFunc functions
func newSeqChecks[T any](initChecks ...checkFunc[T]) *seqChecks[T] {
	return &seqChecks[T]{
		c: newValChecks[T](initChecks...),
	}
}

// Append adds another checkFunc to a seqChecks collection
func (sc *seqChecks[T]) Append(check func(T, *with.ValidationOptions) error) {
	sc.c.Append(check)
}

// Evaluate runs every checkFunc against a value and returns the first error
func (sc *seqChecks[T]) Evaluate(seq iter.Seq[T], opts *with.ValidationOptions) error {
	if opts.CollectAllErrors() {
		vErrs := newValidationErrors()

		for v := range seq {
			if err := sc.c.Evaluate(v, opts); err != nil {
				vErrs.Append(err)
			}
		}

		if vErrs.HasErrors() {
			return vErrs
		}
	} else {
		for v := range seq {
			if err := sc.c.Evaluate(v, opts); err != nil {
				return err
			}
		}
	}

	return nil
}
