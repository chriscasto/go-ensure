package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
)

// iterable describes values that can be treated as sequences
type iterable[K comparable, V any] interface {
	~[]V | ~map[K]V
}

// lengthy describes types that the len() function will accept
type lengthy[K comparable, V any] interface {
	string | iterable[K, V]
}

// checkFunc defines a function that performs a check against a value
type checkFunc[T any] func(T, *with.ValidationOptions) error

// valChecks is a collection of checkFunc functions
type valChecks[T any] struct {
	c []checkFunc[T]
}

// newValChecks creates a new valChecks struct with an optional set of initial checkFunc functions
func newValChecks[T any](initChecks ...checkFunc[T]) *valChecks[T] {
	return &valChecks[T]{
		c: initChecks,
	}
}

// Append adds another checkFunc to a valChecks collection
func (vc *valChecks[T]) Append(check checkFunc[T]) {
	vc.c = append(vc.c, check)
}

// Count returns the number of functions in the collection
func (vc *valChecks[T]) Count() int {
	return len(vc.c)
}

// Evaluate runs every checkFunc against a value and returns the first error
func (vc *valChecks[T]) Evaluate(target T, opts *with.ValidationOptions) error {
	// handle stored lenValChecks first
	for _, fn := range vc.c {
		if err := fn(target, opts); err != nil {
			return err
		}
	}

	return nil
}

// EvaluateAll runs every checkFunc against a value and collects all errors returned in a ValidationErrors struct
func (vc *valChecks[T]) EvaluateAll(target T, opts *with.ValidationOptions) *ValidationErrors {
	vErrs := newValidationErrors()

	// handle stored lenValChecks first
	for _, fn := range vc.c {
		if err := fn(target, opts); err != nil {
			vErrs.Append(err)
		}
	}

	// always return vErrs, even if none are added
	// nil lenValChecks against interfaces in Go are stupidly error-prone, so need this to avoid segfaults
	return vErrs
}

// lenChecks is a wrapper around valChecks that adds a set of checks around the length of a value
type lenChecks[K comparable, V any, T lengthy[K, V]] struct {
	parent    *valChecks[T]
	lenChecks *valChecks[int]
}

// newLenChecks returns a new instance of a lenChecks struct wrapped around a parent valChecks
func newLenChecks[K comparable, V any, T lengthy[K, V]](parent *valChecks[T]) *lenChecks[K, V, T] {
	return &lenChecks[K, V, T]{
		parent:    parent,
		lenChecks: newValChecks[int](),
	}
}

// addLenCheck adds a length check function to the list
func (lc *lenChecks[K, V, T]) addLenCheck(lenCheck func(int, *with.ValidationOptions) error) {
	// if this is the first length check we're adding, also add a val check to the main list to eval the len checks
	if lc.lenChecks.Count() == 0 {
		lc.parent.Append(func(val T, opts *with.ValidationOptions) error {
			if opts.CollectAllErrors() {
				return lc.lenChecks.EvaluateAll(len(val), opts)
			}

			return lc.lenChecks.Evaluate(len(val), opts)
		})
	}

	lc.lenChecks.Append(lenCheck)
}

// AddIsEmpty adds a length check that asserts length is 0
func (lc *lenChecks[K, V, T]) AddIsEmpty() {
	lc.addLenCheck(func(l int, _ *with.ValidationOptions) error {
		if l != 0 {
			return errors.New(`must be empty`)
		}
		return nil
	})
}

// AddIsNotEmpty adds a length check that asserts length is not 0
func (lc *lenChecks[K, V, T]) AddIsNotEmpty() {
	lc.addLenCheck(func(l int, _ *with.ValidationOptions) error {
		if l == 0 {
			return errors.New(`must not be empty`)
		}
		return nil
	})
}

// AddHasLength adds a length check that asserts length is equal to the provided int
func (lc *lenChecks[K, V, T]) AddHasLength(i int) {
	lc.addLenCheck(func(l int, _ *with.ValidationOptions) error {
		if l != i {
			return fmt.Errorf(`length must equal %d; got %d`, i, l)
		}
		return nil
	})
}

// AddIsLongerThan adds a length check that asserts length is greater than provided int
func (lc *lenChecks[K, V, T]) AddIsLongerThan(i int) {
	lc.addLenCheck(func(l int, _ *with.ValidationOptions) error {
		if l <= i {
			return fmt.Errorf(`must have a length greater than %d; got %d`, i, l)
		}
		return nil
	})
}

// AddIsShorterThan adds a length check that asserts length is less than the provided int
func (lc *lenChecks[K, V, T]) AddIsShorterThan(i int) {
	lc.addLenCheck(func(l int, _ *with.ValidationOptions) error {
		if l >= i {
			return fmt.Errorf(`must have a length less than %d; got %d`, i, l)
		}
		return nil
	})
}

// AddHasLengthWhere adds a length check based on the passed int validator
func (lc *lenChecks[K, V, T]) AddHasLengthWhere(nv *NumberValidator[int]) {
	lc.addLenCheck(func(l int, opts *with.ValidationOptions) error {
		if err := nv.Validate(l, opts); err != nil {
			return err
		}
		return nil
	})
}
