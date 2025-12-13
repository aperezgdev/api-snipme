package domain

import (
	"testing"
	"time"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestNewUser(t *testing.T) {
	t.Run("Success on valid parameters", func(t *testing.T) {
		t.Parallel()
		email := "user@example.com"
		oauthProvider := OAuthProviderGoogle
		oauthSubject := "google-subject-123"

		user, err := NewUser(email, oauthProvider, oauthSubject)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(user.Email) != email {
			t.Errorf("Expected Email %s, got %s", email, string(user.Email))
		}
		if user.OAuthProvider != oauthProvider {
			t.Errorf("Expected OAuthProvider %s, got %s", oauthProvider, user.OAuthProvider)
		}
		if string(user.OAuthSubject) != oauthSubject {
			t.Errorf("Expected OAuthSubject %s, got %s", oauthSubject, string(user.OAuthSubject))
		}
		if user.Id.String() == "" {
			t.Error("Expected ID to be generated")
		}
		if time.Time(user.CreatedOn).IsZero() {
			t.Error("Expected CreatedOn to be set")
		}
		events := user.PullDomainEvents()
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
	})

	t.Run("Success with GitHub provider", func(t *testing.T) {
		t.Parallel()
		email := "user@example.com"
		oauthProvider := OAuthProviderGitHub
		oauthSubject := "github-subject-456"

		user, err := NewUser(email, oauthProvider, oauthSubject)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user.OAuthProvider != OAuthProviderGitHub {
			t.Errorf("Expected OAuthProvider %s, got %s", OAuthProviderGitHub, user.OAuthProvider)
		}
	})

	t.Run("Fails on invalid email", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			email    string
			expected string
		}{
			{
				name:     "empty email",
				email:    "",
				expected: "must be a valid email address",
			},
			{
				name:     "invalid email format",
				email:    "not-an-email",
				expected: "must be a valid email address",
			},
			{
				name:     "email without @",
				email:    "nodomain.com",
				expected: "must be a valid email address",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := NewUser(tt.email, OAuthProviderGoogle, "subject")
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

	t.Run("Fails on empty oauth subject", func(t *testing.T) {
		t.Parallel()
		email := "user@example.com"
		oauthProvider := OAuthProviderGoogle
		oauthSubject := ""

		_, err := NewUser(email, oauthProvider, oauthSubject)
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
