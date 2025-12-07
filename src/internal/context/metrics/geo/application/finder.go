package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type CountryFinder struct {
	logger            shared_domain_context.Logger
	countryRepository domain.CountryRepository
}

func NewCountryFinder(
	logger shared_domain_context.Logger,
	countryRepository domain.CountryRepository,
) *CountryFinder {
	return &CountryFinder{
		logger:            logger,
		countryRepository: countryRepository,
	}
}

func (f *CountryFinder) Run(ctx context.Context, ip string) (*domain.Country, error) {
	f.logger.Info(ctx, "CountryFinder - FindByIp - Finding country by IP", shared_domain_context.NewField("ip", ip))

	ipVO, err := shared_domain.NewIp(ip)
	if err != nil {
		f.logger.Error(ctx, "CountryFinder - FindByIp - Invalid IP address", shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("error", err.Error()))
		return &domain.Country{}, err
	}

	countryOpt, err := f.countryRepository.FindByIp(ctx, ipVO)
	if err != nil {
		f.logger.Error(ctx, "CountryFinder - FindByIp - Error finding country by IP", shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("error", err.Error()))
		return &domain.Country{}, err
	}

	if !countryOpt.IsPresent() {
		f.logger.Info(ctx, "CountryFinder - FindByIp - No country found for IP", shared_domain_context.NewField("ip", ip))
		return &domain.Country{}, shared_domain_context.NewNotFoundError("no country found for the given IP")
	}

	country := countryOpt.Get()
	f.logger.Info(ctx, "CountryFinder - FindByIp - Country found for IP", shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("country_iso_code", string(country.ISOCode)))

	return country, nil
}
