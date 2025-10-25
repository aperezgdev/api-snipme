package domain

import (
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type CountryCode string

func NewCountryCode(code string) (CountryCode, error) {
	err := validate(code)
	if err != nil {
		return "", err
	}
	return CountryCode(code), nil
}

func validate(code string) error {
	if len(code) != 2 {
		return shared_domain.NewValidationError("country_code", "Country code must be 2 characters long")
	}
	return nil
}
