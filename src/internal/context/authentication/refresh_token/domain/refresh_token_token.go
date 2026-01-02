package domain

import (
	"strings"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type RefreshTokenToken string

func NewRefreshTokenToken(token string) (RefreshTokenToken, error) {
	if strings.TrimSpace(token) == "" {
		return "", shared_domain.NewValidationError("token", "cannot be empty")
	}
	return RefreshTokenToken(token), nil
}
