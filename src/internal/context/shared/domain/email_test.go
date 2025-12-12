package domain

import "testing"

func TestEmail(t *testing.T) {
	t.Run("NewEmail returns valid email", func(t *testing.T) {
		t.Parallel()
		email, err := NewEmail("test@example.com")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if email != "test@example.com" {
			t.Errorf("Expected email to be 'test@example.com', got %v", email)
		}
	})

	t.Run("NewEmail returns error for invalid email", func(t *testing.T) {
		t.Parallel()
		_, err := NewEmail("invalid-email")
		if err == nil {
			t.Errorf("Expected error for invalid email")
		}
	})
}
