package domain

import (
	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
)

type Country struct {
	ISOCode domain_shared.CountryCode
}

func NewCountry(isoCode string) (Country, error) {
	countryCodeVO, err := domain_shared.NewCountryCode(isoCode)
	if err != nil {
		return Country{}, err
	}
	return Country{
		ISOCode: countryCodeVO,
	}, nil
}
