package domain

import "testing"

func TestErrors(t *testing.T) {
	t.Run("NewValidationError", func(t *testing.T) {
		t.Parallel()
		err := NewValidationError("field", "message")
		if err.Field != "field" || err.Message != "message" {
			t.Errorf("Unexpected ValidationError: %+v", err)
		}
		if err.Error() != "field: message" {
			t.Errorf("Unexpected error string: %s", err.Error())
		}
	})

	t.Run("NewNotFoundError", func(t *testing.T) {
		t.Parallel()
		err := NewNotFoundError("not found")
		if err.Message != "not found" {
			t.Errorf("Unexpected NotFoundError: %+v", err)
		}
		if err.Error() != "Not Found: not found" {
			t.Errorf("Unexpected error string: %s", err.Error())
		}
	})
}
