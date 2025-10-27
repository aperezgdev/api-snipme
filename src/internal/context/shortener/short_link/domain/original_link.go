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
	url, err := url.ParseRequestURI(link)
	if err != nil || url.Scheme == "" || url.Host == "" {
		return shared_domain.NewValidationError("original_link", "must be a valid URL")
	}
	return err
}
