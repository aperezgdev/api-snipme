package domain

import (
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewOAuthProvider(t *testing.T) {
	t.Run("Success on valid providers", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected OAuthProvider
		}{
			{
				name:     "google lowercase",
				input:    "google",
				expected: OAuthProviderGoogle,
			},
			{
				name:     "google uppercase",
				input:    "GOOGLE",
				expected: OAuthProviderGoogle,
			},
			{
				name:     "google mixed case",
				input:    "Google",
				expected: OAuthProviderGoogle,
			},
			{
				name:     "github lowercase",
				input:    "github",
				expected: OAuthProviderGitHub,
			},
			{
				name:     "github uppercase",
				input:    "GITHUB",
				expected: OAuthProviderGitHub,
			},
			{
				name:     "github mixed case",
				input:    "GitHub",
				expected: OAuthProviderGitHub,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				provider, err := NewOAuthProvider(tt.input)
				if err != nil {
					t.Fatalf("Expected no error for %s, got %v", tt.name, err)
				}
				if provider != tt.expected {
					t.Errorf("Expected provider %s, got %s", tt.expected, provider)
				}
			})
		}
	})

	t.Run("Fails on invalid provider", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "empty string",
				input:    "",
				expected: "must be 'google' or 'github'",
			},
			{
				name:     "invalid provider",
				input:    "facebook",
				expected: "must be 'google' or 'github'",
			},
			{
				name:     "twitter",
				input:    "twitter",
				expected: "must be 'google' or 'github'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewOAuthProvider(tt.input)
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

func TestOAuthProvider_String(t *testing.T) {
	t.Run("Returns correct string representation", func(t *testing.T) {
		t.Parallel()
		if OAuthProviderGoogle.String() != "google" {
			t.Errorf("Expected 'google', got %s", OAuthProviderGoogle.String())
		}
		if OAuthProviderGitHub.String() != "github" {
			t.Errorf("Expected 'github', got %s", OAuthProviderGitHub.String())
		}
	})
}
