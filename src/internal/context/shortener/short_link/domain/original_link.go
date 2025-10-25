package domain

import (
	"net/url"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type ShortLinkOriginalRoute string

func NewShortLinkOriginalRoute(link string) (ShortLinkOriginalRoute, error) {
	if err := validate(link); err != nil {
		return "", err
	}
	return ShortLinkOriginalRoute(link), nil
}

func validate(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return shared_domain.NewValidationError("original_link", "must be a valid URL")
	}
	return err
}
