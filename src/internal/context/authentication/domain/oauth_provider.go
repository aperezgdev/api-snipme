package domain

import (
	"strings"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderGitHub OAuthProvider = "github"
)

func NewOAuthProvider(provider string) (OAuthProvider, error) {
	normalized := OAuthProvider(strings.ToLower(provider))
	
	switch normalized {
	case OAuthProviderGoogle, OAuthProviderGitHub:
		return normalized, nil
	default:
		return "", shared_domain.NewValidationError("oauth_provider", "must be 'google' or 'github'")
	}
}

func (p OAuthProvider) String() string {
	return string(p)
}
