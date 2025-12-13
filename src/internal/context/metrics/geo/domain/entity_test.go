package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewCountry(t *testing.T) {
	t.Run("Success on valid ISO code", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			isoCode string
		}{
			{
				name:    "US code",
				isoCode: "US",
			},
			{
				name:    "GB code",
				isoCode: "GB",
			},
			{
				name:    "ES code",
				isoCode: "ES",
			},
			{
				name:    "lowercase code",
				isoCode: "de",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				country, err := NewCountry(tt.isoCode)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(country.ISOCode) != tt.isoCode {
					t.Errorf("Expected ISO code %s, got %s", tt.isoCode, string(country.ISOCode))
				}
			})
		}
	})

	t.Run("Fails on invalid ISO code", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty code",
				input:    "",
				expected: "Country code must be 2 characters long",
			},
			{
				name:     "single character",
				input:    "U",
				expected: "Country code must be 2 characters long",
			},
			{
				name:     "three characters",
				input:    "USA",
				expected: "Country code must be 2 characters long",
			},
			{
				name:     "too long",
				input:    "USAA",
				expected: "Country code must be 2 characters long",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewCountry(tt.input)
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
