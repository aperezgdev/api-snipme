package domain

import (
	"testing"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewClientName(t *testing.T) {
	t.Run("Success on valid name", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "single word",
				input: "John",
			},
			{
				name:  "multiple words",
				input: "John Doe",
			},
			{
				name:  "with numbers",
				input: "John123",
			},
			{
				name:  "max length",
				input: "12345678901234567890123456789012345678901234567890",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				clientName, err := NewClientName(tt.input)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(clientName) != tt.input {
					t.Errorf("Expected name %s, got %s", tt.input, string(clientName))
				}
			})
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
				_, err := NewClientName(tt.input)
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

func TestValidateName(t *testing.T) {
	t.Run("Returns nil on valid name", func(t *testing.T) {
		t.Parallel()
		err := validateName("Valid Name")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Returns error on empty name", func(t *testing.T) {
		t.Parallel()
		err := validateName("")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("Returns error on name too long", func(t *testing.T) {
		t.Parallel()
		err := validateName("This is a very long name that exceeds the maximum allowed length of fifty characters")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
