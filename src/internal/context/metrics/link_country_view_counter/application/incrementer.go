package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkCountryViewCounterIncrementer struct {
	logger   shared_domain_context.Logger
	repo     domain.LinkCountryViewCounterRepository
	finder   domain.LinkCountryViewCounterFinder
	eventBus shared_domain_context.EventBus
}

func NewLinkCountryViewCounterIncrementer(logger shared_domain_context.Logger, repo domain.LinkCountryViewCounterRepository, eventBus shared_domain_context.EventBus) *LinkCountryViewCounterIncrementer {
	return &LinkCountryViewCounterIncrementer{
		logger:   logger,
		repo:     repo,
		finder:   domain.NewLinkCountryViewCounterFinder(logger, repo),
		eventBus: eventBus,
	}
}

func (u LinkCountryViewCounterIncrementer) Run(ctx context.Context, linkId, countryCode string, uniqueViews, totalViews uint) (*domain.LinkCountryViewCounter, error) {
	u.logger.Info(ctx, "LinkCountryViewCounterUpdater - Run - Params into: ", shared_domain_context.NewField("linkCountryViewCounter", linkId), shared_domain_context.NewField("uniqueViews", uniqueViews), shared_domain_context.NewField("totalViews", totalViews))

	linkCountryViewCounterOptional, err := u.finder.Run(ctx, linkId, countryCode)

	if !linkCountryViewCounterOptional.IsPresent() {
		u.logger.Info(ctx, "LinkCountryViewCounterUpdater - Run - Link country view counter not found, creating new one", shared_domain_context.NewField("linkCountryViewCounter", linkId))
		return nil, shared_domain_context.NewNotFoundError("Link country view counter not found")
	}

	linkCountryViewCounter := linkCountryViewCounterOptional.Get()
	linkCountryViewCounter.Increment(totalViews, uniqueViews)
	u.logger.Info(ctx, "LinkCountryViewCounterUpdater - Run - Incremented link country view counter", shared_domain_context.NewField("linkCountryViewCounter", linkId), shared_domain_context.NewField("uniqueViews", linkCountryViewCounter.UniqueViews), shared_domain_context.NewField("totalViews", linkCountryViewCounter.TotalViews))

	err = u.repo.Update(ctx, linkCountryViewCounter)
	if err != nil {
		u.logger.Error(ctx, "LinkCountryViewCounterUpdater - Run - Error updating link country view counter", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	u.logger.Info(ctx, "LinkCountryViewCounterUpdater - Run - Successfully updated link country view counter", shared_domain_context.NewField("linkCountryViewCounter", linkId))

	u.eventBus.Publish(ctx, linkCountryViewCounter.PullDomainEvents()...)

	return &linkCountryViewCounter, nil
}
