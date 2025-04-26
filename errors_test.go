package ensure_test

import (
	"errors"
	"github.com/chriscasto/go-ensure"
	"strings"
	"testing"
)

// TestValidationErrors_IsError checks to make sure ValidatorErrors implements the error interface
func TestValidationErrors_IsError(t *testing.T) {
	var _ error = &ensure.ValidationErrors{}
}

func TestErrorAsValidationErrors(t *testing.T) {
	t.Run("is validation errors", func(t *testing.T) {
		retValErr := func() error {
			return ensure.NewValidationErrors()
		}

		err := retValErr()
		vErrs := ensure.ErrorAsValidationErrors(err)

		if vErrs == nil {
			t.Errorf("expected validation errors; got none")
		}
	})

	t.Run("is not validation errors", func(t *testing.T) {
		retValErr := func() error {
			return errors.New("plain error")
		}

		err := retValErr()
		vErrs := ensure.ErrorAsValidationErrors(err)

		if vErrs != nil {
			t.Errorf("expected no error; got %s", vErrs.Error())
		}
	})
}

func TestValidationErrors_Errors(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		vErrs := ensure.NewValidationErrors()
		expect := "there were no validation errors"

		if vErrs.Error() != expect {
			t.Errorf("expected %s; got %s", expect, vErrs.Error())
		}
	})

	t.Run("type error", func(t *testing.T) {
		vErrs := ensure.NewValidationErrors()
		msg := "this should be anonymized"
		vErrs.Append(ensure.NewTypeError(msg))

		// type errors are anonymized to avoid leaks
		if !strings.Contains(vErrs.Error(), "type error") {
			t.Errorf(`expected error to contain "type error"; got "%s"`, vErrs.Error())
		}
	})

	t.Run("validation error", func(t *testing.T) {
		vErrs := ensure.NewValidationErrors()
		msg := "validation error"
		vErrs.Append(ensure.NewValidationError(msg))

		if vErrs.Error() != msg {
			t.Errorf("expected %s; got %s", msg, vErrs.Error())
		}
	})
}

func TestValidationErrors_Append(t *testing.T) {
	// most Append method statements are covered in the Extend tests below
	vErr1 := ensure.NewValidationErrors()
	vErr1.Append(nil)

	if vErr1.HasErrors() {
		t.Errorf(`expected no errors; got some`)
	}

	vErr2 := ensure.NewValidationErrors()
	vErr2.Append(errors.New("normal err 1"))

	if !vErr2.HasValidationErrors() {
		t.Errorf(`expected validation errors; got none`)
	}

	if vErr2.HasTypeErrors() {
		t.Errorf(`expected no type errors; got some`)
	}

	vErr3 := ensure.NewValidationErrors()
	vErr3.Append(ensure.NewTypeError("type err 1"))

	if vErr3.HasValidationErrors() {
		t.Errorf(`expected no validation errors; got some`)
	}

	if !vErr3.HasTypeErrors() {
		t.Errorf(`expected type errors; got none`)
	}
}

func TestValidationErrors_Extend(t *testing.T) {
	nilErr := ensure.NewValidationErrors()
	nilErr.Extend(nil)

	if nilErr.HasErrors() {
		t.Errorf(`expected no errors`)
	}

	vErr1 := ensure.NewValidationErrors()
	vErr1.Append(errors.New("normal err 1"))
	vErr1.Append(errors.New("normal err 2"))
	vErr1.Append(ensure.NewTypeError("type err 1"))

	vErr2 := ensure.NewValidationErrors()
	vErr2.Append(errors.New("normal err 3"))
	vErr2.Append(errors.New("normal err 4"))
	vErr2.Append(ensure.NewTypeError("type err 2"))

	vErr1.Extend(vErr2)

	if !vErr1.HasErrors() {
		t.Errorf(`expected errors; got none`)
	}

	if !vErr1.HasTypeErrors() {
		t.Errorf(`expected type errors; got none`)
	}

	if !vErr1.HasValidationErrors() {
		t.Errorf(`expected validation errors; got none`)
	}

	tErrs := vErr1.TypeErrors()

	if len(tErrs) != 2 {
		t.Errorf(`expected type errors to have length 2; got "%d"`, len(tErrs))
	}

	vErrs := vErr1.ValidationErrors()

	if len(vErrs) != 4 {
		t.Errorf(`expected validation errors to have length 4; got "%d"`, len(tErrs))
	}
}

func TestValidationErrors_GetTypeErrors(t *testing.T) {
	t.Run("type error", func(t *testing.T) {
		vErrs := ensure.NewValidationErrors()
		msg := "type error"
		vErrs.Append(ensure.NewTypeError(msg))

		tErrs := vErrs.TypeErrors()

		if len(tErrs) != 1 {
			t.Errorf(`expected type errors to have length 1; got "%d"`, len(tErrs))
		}

		if tErrs[0].Error() != msg {
			t.Errorf(`expected "%s"; got "%s"`, msg, tErrs[0].Error())
		}
	})
}
