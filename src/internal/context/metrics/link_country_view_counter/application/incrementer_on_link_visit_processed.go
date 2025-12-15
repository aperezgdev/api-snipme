package application

import (
	"context"
	"errors"

	geo_application "github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/application"
	geo_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	domain_link_visit "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type IncremeterOnLinkVisitProcessed struct {
	logger        shared_domain_context.Logger
	countryFinder geo_application.CountryFinder
	incrementer   domain.LinkCountryViewCounterIncrementer
}

func NewIncremeterOnLinkVisitProcessed(
	logger shared_domain_context.Logger,
	repoGeo geo_domain.CountryRepository,
	repo domain.LinkCountryViewCounterRepository,
) *IncremeterOnLinkVisitProcessed {
	return &IncremeterOnLinkVisitProcessed{
		logger:        logger,
		countryFinder: *geo_application.NewCountryFinder(logger, repoGeo),
		incrementer:   *domain.NewLinkCountryViewCounterIncrementer(logger, repo),
	}
}

func (u *IncremeterOnLinkVisitProcessed) On(ctx context.Context, event shared_domain_context.DomainEvent) error {
	u.logger.Info(ctx, "UpserterOnLinkVisitProcessed - On - Received LinkVisitProcessed event", shared_domain_context.NewField("aggregateID", event.AggregateID()))

	eventData, ok := event.(domain_link_visit.LinkVisitsProcessed)
	if !ok {
		u.logger.Error(ctx, "UpserterOnLinkVisitProcessed - On - Invalid event data type")
		return shared_domain_context.NewValidationError("event data", "invalid type")
	}
	u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Processing LinkVisitProcessed event", shared_domain_context.NewField("aggregateID", eventData.AggregateID()), shared_domain_context.NewField("ipsFrequency", eventData.IpsFrequency))

	cacheIp := make(map[string]string)
	reduceViewsByCountry := make(map[string]struct {
		uniqueViews uint
		totalViews  uint
	})

	for _, ip := range eventData.IpsFrequency {
		isoCode := "UNK"
		if cacheIp[ip.Ip] == "" {
			u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Country not in cache for IP, looking up", shared_domain_context.NewField("ip", ip.Ip))
			country, err := u.countryFinder.Run(ctx, ip.Ip)
			if err == nil {
				u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Country found for IP", shared_domain_context.NewField("ip", ip.Ip), shared_domain_context.NewField("country", country.ISOCode))
				isoCode = string(country.ISOCode)
				cacheIp[ip.Ip] = isoCode
			}

			if errors.Is(err, &shared_domain_context.NotFoundError{}) {
				u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Country not found for IP, setting as UNK", shared_domain_context.NewField("ip", ip.Ip))
			} else if err != nil {
				u.logger.Error(ctx, "IncremeterOnLinkVisitProcessed - On - Error finding country for IP", shared_domain_context.NewField("error", err.Error()), shared_domain_context.NewField("ip", ip.Ip))
				return err
			}
		} else {
			u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Country found in cache for IP", shared_domain_context.NewField("ip", ip.Ip))
			isoCode = cacheIp[ip.Ip]
		}

		if viewsByCountry, exists := reduceViewsByCountry[isoCode]; exists {
			u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Aggregating views for country", shared_domain_context.NewField("country", isoCode))
			viewsByCountry.totalViews += ip.TotalViews
			viewsByCountry.uniqueViews += ip.UniqueViews
			reduceViewsByCountry[isoCode] = viewsByCountry
		} else {
			u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Setting initial views for country", shared_domain_context.NewField("country", isoCode))
			reduceViewsByCountry[isoCode] = struct {
				uniqueViews uint
				totalViews  uint
			}{
				uniqueViews: ip.UniqueViews,
				totalViews:  ip.TotalViews,
			}
		}
	}

	for countryCode, views := range reduceViewsByCountry {
		_, err := u.incrementer.Run(ctx, eventData.AggregateID(), countryCode, views.uniqueViews, views.totalViews)
		if err != nil {
			u.logger.Error(ctx, "IncremeterOnLinkVisitProcessed - On - Error incrementing link country view counter", shared_domain_context.NewField("error", err.Error()), shared_domain_context.NewField("linkId", eventData.AggregateID()), shared_domain_context.NewField("countryCode", countryCode))
			return err
		}
		u.logger.Info(ctx, "IncremeterOnLinkVisitProcessed - On - Successfully incremented link country view counter", shared_domain_context.NewField("linkId", eventData.AggregateID()), shared_domain_context.NewField("countryCode", countryCode), shared_domain_context.NewField("uniqueViews", views.uniqueViews), shared_domain_context.NewField("totalViews", views.totalViews))
	}

	return nil
}
