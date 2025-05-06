package with_test

import (
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

func TestAnyOptions_DefaultError(t *testing.T) {
	opts := with.DefaultAnyOptions()
	orig := opts.DefaultError()

	if orig.Error() != with.DefaultAnyValidatorErrorMsg {
		t.Errorf("expected default any error message, got %s", orig.Error())
	}

	errMsg := "this is the new default error message"

	with.AnyOptionDefaultError(errMsg)(opts)

	newErr := opts.DefaultError()

	if newErr.Error() != errMsg {
		t.Errorf("expected new error message, got %s", newErr.Error())
	}
}

func TestAnyOptions_PassThroughErrorsFrom(t *testing.T) {
	opts := with.DefaultAnyOptions()

	ints := []int{1, 2, 3}
	passthru := []int{1, 3}
	expectAllow := map[int]bool{1: true, 2: false, 3: true}

	for _, i := range ints {
		if opts.PassThroughErrorsFrom(i) {
			t.Errorf("expected no values to permit passthough, was permitted for %d", i)
		}
	}

	with.AnyOptionPassThroughErrorsFrom(passthru...)(opts)

	for _, i := range ints {
		if opts.PassThroughErrorsFrom(i) && !expectAllow[i] {
			t.Errorf("expected value [%d] to not permit passthough but was permitted", i)
		} else if !opts.PassThroughErrorsFrom(i) && expectAllow[i] {
			t.Errorf("expected value [%d] to permit passthough but was not permitted", i)
		}
	}
}
