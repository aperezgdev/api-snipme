package domain

import (
	"testing"
	"time"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewClient(t *testing.T) {
	t.Run("Success on valid name and email", func(t *testing.T) {
		t.Parallel()
		name := "John Doe"
		email := "john.doe@example.com"

		client, err := NewClient(name, email)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(client.Name) != name {
			t.Errorf("Expected name %s, got %s", name, string(client.Name))
		}
		if string(client.Email) != email {
			t.Errorf("Expected email %s, got %s", email, string(client.Email))
		}
		if client.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if time.Time(client.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
	})

	t.Run("Fails on invalid name", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty name",
				input:    "",
				expected: "Name cannot be empty",
			},
			{
				name:     "name too long",
				input:    "This is a very long name that exceeds the maximum allowed length of fifty characters",
				expected: "Name cannot be longer than 50 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewClient(tt.input, "valid@example.com")
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
				validationErr, ok := err.(domain_shared.ValidationError)
				if !ok {
					t.Fatalf("Expected ValidationError, got %T", err)
				}
				if validationErr.Message != tt.expected {
					t.Errorf("Expected error message '%s', got '%s'", tt.expected, validationErr.Message)
				}
			})
		}
	})

	t.Run("Fails on invalid email", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty email",
				input:    "",
				expected: "must be a valid email address",
			},
			{
				name:     "invalid email format",
				input:    "not-an-email",
				expected: "must be a valid email address",
			},
			{
				name:     "email without @",
				input:    "nodomain.com",
				expected: "must be a valid email address",
			},
			{
				name:     "email without domain",
				input:    "user@",
				expected: "must be a valid email address",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewClient("Valid Name", tt.input)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
				validationErr, ok := err.(domain_shared.ValidationError)
				if !ok {
					t.Fatalf("Expected ValidationError, got %T", err)
				}
				if validationErr.Message != tt.expected {
					t.Errorf("Expected error message '%s', got '%s'", tt.expected, validationErr.Message)
				}
			})
		}
	})
}
