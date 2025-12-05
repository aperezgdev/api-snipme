package domain

import shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"

type OAuthSubject string

func NewOAuthSubject(subject string) (OAuthSubject, error) {
	if subject == "" {
		return "", shared_domain.NewValidationError("oauth_subject", "cannot be empty")
	}
	return OAuthSubject(subject), nil
}
