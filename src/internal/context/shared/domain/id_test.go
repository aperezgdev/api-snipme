package domain

import (
	"testing"
)

func TestID(t *testing.T) {
	t.Run("NewID returns valid ID", func(t *testing.T) {
		t.Parallel()
		id, err := NewID()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if id.String() == "" {
			t.Errorf("Expected non-empty string for ID")
		}
	})

	t.Run("ParseID parses valid and invalid strings", func(t *testing.T) {
		t.Parallel()
		id, _ := NewID()
		parsed, err := ParseID(id.String())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if parsed != id {
			t.Errorf("Expected parsed ID to equal original")
		}

		_, err = ParseID("invalid-uuid")
		if err == nil {
			t.Errorf("Expected error for invalid uuid string")
		}
	})
}
