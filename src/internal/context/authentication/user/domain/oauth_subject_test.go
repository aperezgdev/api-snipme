package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewOAuthSubject(t *testing.T) {
	t.Run("Success on valid subject", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			subject string
		}{
			{
				name:    "simple string",
				subject: "user-123",
			},
			{
				name:    "google subject format",
				subject: "google-oauth2|123456789",
			},
			{
				name:    "github subject format",
				subject: "github|username",
			},
			{
				name:    "with special characters",
				subject: "subject-with_special.chars@123",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				subject, err := NewOAuthSubject(tt.subject)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(subject) != tt.subject {
					t.Errorf("Expected subject %s, got %s", tt.subject, string(subject))
				}
			})
		}
	})

	t.Run("Fails on empty subject", func(t *testing.T) {
		t.Parallel()
		_, err := NewOAuthSubject("")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		validationErr, ok := err.(shared_domain.ValidationError)
		if !ok {
			t.Fatalf("Expected ValidationError, got %T", err)
		}
		if validationErr.Message != "cannot be empty" {
			t.Errorf("Expected error message 'cannot be empty', got '%s'", validationErr.Message)
		}
	})
}
