package domain

import (
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type ShortLinkSummary string

func NewShortLinkSummary(summary string) (ShortLinkSummary, error) {
	if len(summary) > 255 {
		return "", shared_domain_context.NewValidationError("summary", "must not exceed 255 characters")
	}
	return ShortLinkSummary(summary), nil
}
