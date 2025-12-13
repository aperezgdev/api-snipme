package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewRefreshTokenToken(t *testing.T) {
	t.Run("Success on valid token", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			token string
		}{
			{
				name:  "simple token",
				token: "token-123",
			},
			{
				name:  "JWT format",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			},
			{
				name:  "long token",
				token: "very-long-token-with-many-characters-1234567890-abcdefghijklmnopqrstuvwxyz",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				token, err := NewRefreshTokenToken(tt.token)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(token) != tt.token {
					t.Errorf("Expected token %s, got %s", tt.token, string(token))
				}
			})
		}
	})

	t.Run("Fails on empty or whitespace token", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			token    string
			expected string
		}{
			{
				name:     "empty string",
				token:    "",
				expected: "cannot be empty",
			},
			{
				name:     "only spaces",
				token:    "   ",
				expected: "cannot be empty",
			},
			{
				name:     "only tabs",
				token:    "\t\t",
				expected: "cannot be empty",
			},
			{
				name:     "mixed whitespace",
				token:    " \t \n ",
				expected: "cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewRefreshTokenToken(tt.token)
				if err == nil {
					t.Fatalf("Expected error for %s, got nil", tt.name)
				}
				validationErr, ok := err.(shared_domain.ValidationError)
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
