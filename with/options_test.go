package with_test

import (
	"github.com/chriscasto/go-ensure/with"
	"testing"
)

func TestValidationOptions_CollectAllErrors(t *testing.T) {
	defOpts := with.ValidationOptions{}

	if defOpts.CollectAllErrors() {
		t.Errorf("expected default options to not collect all errors")
	}

	opts := with.Options(
		with.OptionCollectAllErrors(),
	)

	if !opts.CollectAllErrors() {
		t.Errorf("expected OptionCollectAllErrors to result in collecting all errors")
	}
}
