package domain

import (
	"testing"
	"time"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewShortLink(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		summary := "Test Summary"
		originalLink := "https://example.com/path"
		client := "00000000-0000-0000-0000-000000000000"

		shortLink, err := NewShortLink(summary, originalLink, client)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(shortLink.Summary) != summary {
			t.Errorf("Expected summary %s, got %s", summary, string(shortLink.Summary))
		}
		if string(shortLink.OriginalRoute) != originalLink {
			t.Errorf("Expected original route %s, got %s", originalLink, string(shortLink.OriginalRoute))
		}
		if shortLink.Client.String() != client {
			t.Errorf("Expected client %s, got %s", client, shortLink.Client.String())
		}
		if shortLink.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if string(shortLink.Code) == "" {
			t.Error("Expected Code to be generated")
		}
		if time.Time(shortLink.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
		events := shortLink.PullDomainEvents()
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
	})

	t.Run("Success without client", func(t *testing.T) {
		t.Parallel()
		summary := "Test Summary"
		originalLink := "https://example.com/path"
		client := ""

		shortLink, err := NewShortLink(summary, originalLink, client)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if shortLink.Client.String() != "00000000-0000-0000-0000-000000000000" {
			t.Errorf("Expected empty client UUID, got %s", shortLink.Client.String())
		}
	})

	t.Run("Fails on invalid summary", func(t *testing.T) {
		t.Parallel()
		summary := "This is a very long summary that exceeds the maximum allowed length of 255 characters. This is a very long summary that exceeds the maximum allowed length of 255 characters. This is a very long summary that exceeds the maximum allowed length of 255 characters. This is a very long summary that exceeds the maximum allowed length."
		originalLink := "https://example.com/path"
		client := ""

		_, err := NewShortLink(summary, originalLink, client)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		validationErr, ok := err.(shared_domain.ValidationError)
		if !ok {
			t.Fatalf("Expected ValidationError, got %T", err)
		}
		if validationErr.Message != "must not exceed 255 characters" {
			t.Errorf("Expected error message 'must not exceed 255 characters', got '%s'", validationErr.Message)
		}
	})

	t.Run("Fails on invalid original link", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty link",
				input:    "",
				expected: "must be a valid URL",
			},
			{
				name:     "invalid URL",
				input:    "not-a-url",
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
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewShortLink("Valid Summary", tt.input, "")
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

	t.Run("Fails on invalid client ID", func(t *testing.T) {
		t.Parallel()
		summary := "Test Summary"
		originalLink := "https://example.com/path"
		client := "invalid-uuid"

		_, err := NewShortLink(summary, originalLink, client)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestNewPublicShortLink(t *testing.T) {
	t.Run("Success on valid original link", func(t *testing.T) {
		t.Parallel()
		originalLink := "https://example.com/path"

		shortLink, err := NewPublicShortLink(originalLink)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(shortLink.OriginalRoute) != originalLink {
			t.Errorf("Expected original route %s, got %s", originalLink, string(shortLink.OriginalRoute))
		}
		if shortLink.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if string(shortLink.Code) == "" {
			t.Error("Expected Code to be generated")
		}
		if time.Time(shortLink.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
		events := shortLink.PullDomainEvents()
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
	})

	t.Run("Fails on invalid original link", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty link",
				input:    "",
				expected: "must be a valid URL",
			},
			{
				name:     "invalid URL",
				input:    "not-a-url",
				expected: "must be a valid URL",
			},
			{
				name:     "URL without scheme",
				input:    "example.com",
				expected: "must be a valid URL",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewPublicShortLink(tt.input)
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
