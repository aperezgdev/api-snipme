package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewShortLinkOriginalRoute(t *testing.T) {
	t.Run("Success on valid URL", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "HTTP URL",
				input: "http://example.com",
			},
			{
				name:  "HTTPS URL",
				input: "https://example.com",
			},
			{
				name:  "URL with path",
				input: "https://example.com/path/to/resource",
			},
			{
				name:  "URL with query parameters",
				input: "https://example.com?param=value",
			},
			{
				name:  "URL with port",
				input: "https://example.com:8080/path",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				route, err := NewShortLinkOriginalRoute(tt.input)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(route) != tt.input {
					t.Errorf("Expected route %s, got %s", tt.input, string(route))
				}
			})
		}
	})

	t.Run("Fails on invalid URL", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty URL",
				input:    "",
				expected: "must be a valid URL",
			},
			{
				name:     "URL without scheme",
				input:    "example.com",
				expected: "must be a valid URL",
			},
			{
				name:     "URL without host",
				input:    "https://",
				expected: "must be a valid URL",
			},
			{
				name:     "invalid URL format",
				input:    "not a url",
				expected: "must be a valid URL",
			},
			{
				name:     "URL with only scheme",
				input:    "http://",
				expected: "must be a valid URL",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewShortLinkOriginalRoute(tt.input)
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

func TestValidate(t *testing.T) {
	t.Run("Returns nil on valid URL", func(t *testing.T) {
		t.Parallel()
		err := validate("https://example.com")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Returns error on invalid URL", func(t *testing.T) {
		t.Parallel()
		err := validate("not-a-url")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("Returns error on empty URL", func(t *testing.T) {
		t.Parallel()
		err := validate("")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
