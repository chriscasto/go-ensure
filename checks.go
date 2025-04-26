package ensure

import (
	"errors"
	"fmt"
	"github.com/chriscasto/go-ensure/with"
	"iter"
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

// Evaluate runs every checkFunc against a value and returns any errors
func (vc *valChecks[T]) Evaluate(target T, opts *with.ValidationOptions) error {
	if opts.CollectAllErrors() {
		vErrs := newValidationErrors()

		for _, fn := range vc.c {
			if err := fn(target, opts); err != nil {
				vErrs.Append(err)
			}
		}

		if vErrs.HasErrors() {
			return vErrs
		}

		return nil
	} else {
		for _, fn := range vc.c {
			if err := fn(target, opts); err != nil {
				return err
			}
		}
	}

	return nil
}

// lenChecks is an extension of valChecks that can set additional checks around the length of a value
type lenChecks[K comparable, V any, T lengthy[K, V]] struct {
	*valChecks[T]
	lenChecks *valChecks[int]
}

// newLenChecks returns a new instance of a lenChecks struct wrapped around a parent valChecks
func newLenChecks[K comparable, V any, T lengthy[K, V]]() *lenChecks[K, V, T] {
	return &lenChecks[K, V, T]{
		valChecks: newValChecks[T](),
		lenChecks: newValChecks[int](),
	}
}

// addLenCheck adds a length check function to the list
func (lc *lenChecks[K, V, T]) addLenCheck(lenCheck func(int, *with.ValidationOptions) error) {
	// if this is the first length check we're adding, also add a val check to the main list to eval the len checks
	if lc.lenChecks.Count() == 0 {
		lc.Append(func(val T, opts *with.ValidationOptions) error {
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

// iterChecks is an extension of lenChecks that can evaluate checks against each key/value in the sequence it contains
type iterChecks[K comparable, V any, T iterable[K, V]] struct {
	*lenChecks[K, V, T]
	iterKeyChecks *valChecks[K]
	iterValChecks *valChecks[V]
	toSeq         func(T) iter.Seq2[K, V]
	checkAdded    bool
}

// newIterChecks creates a new instance of iterChecks
func newIterChecks[K comparable, V any, T iterable[K, V]](toSeq func(T) iter.Seq2[K, V]) *iterChecks[K, V, T] {
	return &iterChecks[K, V, T]{
		lenChecks:     newLenChecks[K, V, T](),
		iterKeyChecks: newValChecks[K](),
		iterValChecks: newValChecks[V](),
		toSeq:         toSeq,
	}
}

// newArrIterChecks creates a new instance of iterChecks specific to an array
func newArrIterChecks[V any]() *iterChecks[int, V, []V] {
	return newIterChecks[int, V, []V](
		func(arr []V) iter.Seq2[int, V] {
			return func(yield func(int, V) bool) {
				for k, v := range arr {
					if !yield(k, v) {
						return
					}
				}
			}
		})
}

// newMapIterChecks creates a new instance of iterChecks specific to a map
func newMapIterChecks[K comparable, V any]() *iterChecks[K, V, map[K]V] {
	return newIterChecks[K, V, map[K]V](
		func(it map[K]V) iter.Seq2[K, V] {
			return func(yield func(K, V) bool) {
				for k, v := range it {
					if !yield(k, v) {
						return
					}
				}
			}
		})
}

// addIterSeqCheck appends a check to the root valCheck that will evaluate each individual value in the sequence
func (ic *iterChecks[K, V, T]) addIterSeqCheck() {
	// don't add another check if we've already added one
	if ic.checkAdded {
		return
	}

	// append the check that will evaluate keys and values to the main list of checks
	ic.Append(func(it T, opts *with.ValidationOptions) error {
		seq := ic.toSeq(it)

		if opts.CollectAllErrors() {
			vErrs := newValidationErrors()

			for k, v := range seq {
				if err := ic.iterKeyChecks.Evaluate(k, opts); err != nil {
					vErrs.Append(err)
				}

				if err := ic.iterValChecks.Evaluate(v, opts); err != nil {
					vErrs.Append(err)
				}
			}

			if vErrs.HasErrors() {
				return vErrs
			}
		} else {
			for k, v := range seq {
				if err := ic.iterKeyChecks.Evaluate(k, opts); err != nil {
					return err
				}

				if err := ic.iterValChecks.Evaluate(v, opts); err != nil {
					return err
				}
			}
		}

		return nil
	})

	// make sure to mark that the check has been added so it doesn't get added again
	ic.checkAdded = true
}

// AddIterKeyCheck adds a check against the keys of this iterable
func (ic *iterChecks[K, V, T]) AddIterKeyCheck(check func(K, *with.ValidationOptions) error) {
	ic.addIterSeqCheck()
	ic.iterKeyChecks.Append(check)
}

// AddIterKeyValidator adds a check to evaluate a validator against the iterable's keys
func (ic *iterChecks[K, V, T]) AddIterKeyValidator(v with.Validator[K]) {
	ic.AddIterKeyCheck(func(val K, opts *with.ValidationOptions) error {
		return v.Validate(val, opts)
	})
}

// AddIterValCheck adds a check against the values of this iterable
func (ic *iterChecks[K, V, T]) AddIterValCheck(check func(V, *with.ValidationOptions) error) {
	ic.addIterSeqCheck()
	ic.iterValChecks.Append(check)
}

// AddIterValValidator adds a check to evaluate a validator against the iterable's values
func (ic *iterChecks[K, V, T]) AddIterValValidator(v with.Validator[V]) {
	ic.AddIterValCheck(func(val V, opts *with.ValidationOptions) error {
		return v.Validate(val, opts)
	})
}
