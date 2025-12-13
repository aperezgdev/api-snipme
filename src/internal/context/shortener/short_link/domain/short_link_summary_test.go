package domain

import (
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewShortLinkSummary(t *testing.T) {
	t.Run("Success on valid summary", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "empty summary",
				input: "",
			},
			{
				name:  "short summary",
				input: "Short summary",
			},
			{
				name:  "max length summary",
				input: "This is a summary with exactly 255 characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut.",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				summary, err := NewShortLinkSummary(tt.input)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if string(summary) != tt.input {
					t.Errorf("Expected summary %s, got %s", tt.input, string(summary))
				}
			})
		}
	})

	t.Run("Fails on summary too long", func(t *testing.T) {
		t.Parallel()
		summary := "This is a summary with more than 255 characters. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip."

		_, err := NewShortLinkSummary(summary)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		validationErr, ok := err.(shared_domain_context.ValidationError)
		if !ok {
			t.Fatalf("Expected ValidationError, got %T", err)
		}
		if validationErr.Message != "must not exceed 255 characters" {
			t.Errorf("Expected error message 'must not exceed 255 characters', got '%s'", validationErr.Message)
		}
	})
}
