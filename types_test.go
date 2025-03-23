package ensure_test

import (
	"github.com/chriscasto/go-ensure"
	"testing"
)

func TestTypeError_Error(t *testing.T) {
	// This is really just a formality to quiet the coverage checker
	msg := "test"
	err := ensure.NewTypeError(msg)

	if err.Error() != msg {
		t.Errorf("TypeError.Error() = `%s`, want `%s`", err.Error(), msg)
	}
}
