package domain

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkCountryViewCounterIncrementer struct {
	logger shared_domain_context.Logger
	repo   LinkCountryViewCounterRepository
	finder LinkCountryViewCounterFinder
}

func NewLinkCountryViewCounterIncrementer(logger shared_domain_context.Logger, repo LinkCountryViewCounterRepository) *LinkCountryViewCounterIncrementer {
	return &LinkCountryViewCounterIncrementer{
		logger: logger,
		repo:   repo,
		finder: NewLinkCountryViewCounterFinder(logger, repo),
	}
}

func (u LinkCountryViewCounterIncrementer) Run(ctx context.Context, linkId, countryCode string, uniqueViews, totalViews uint) (*LinkCountryViewCounter, error) {
	u.logger.Info(ctx, "LinkCountryViewCounterIncrementer - Run - Params into: ", shared_domain_context.NewField("linkCountryViewCounter", linkId), shared_domain_context.NewField("uniqueViews", uniqueViews), shared_domain_context.NewField("totalViews", totalViews))

	linkCountryViewCounterOptional, err := u.finder.Run(ctx, linkId, countryCode)
	if err != nil {
		u.logger.Error(ctx, "LinkCountryViewCLinkCountryViewCounterIncrementerounterUpdater - Run - Error finding link country view counter", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	if !linkCountryViewCounterOptional.IsPresent() {
		u.logger.Info(ctx, "LinkCountryViewCounterIncrementer - Run - Link country view counter not found, creating new one", shared_domain_context.NewField("linkCountryViewCounter", linkId))
		
		newLinkCountryViewCounter, err  := NewLinkCountryViewCounter(linkId, countryCode)
		if err != nil {
			u.logger.Error(ctx, "LinkCountryViewCounterIncrementer - Run - Error creating new link country view counter", shared_domain_context.NewField("error", err.Error()))
			return nil,  err
		}

		newLinkCountryViewCounter.Increment(totalViews, uniqueViews)
		err = u.repo.Save(ctx, *newLinkCountryViewCounter)
		if  err != nil {
			u.logger.Error(ctx, "LinkCountryViewCounterIncrementer - Run - Error saving new link country view counter", shared_domain_context.NewField("error", err.Error()))
			return nil, err
		}

		u.logger.Info(ctx, "LinkCountryViewCounterIncrementer - Run - Successfully created new link country view counter", shared_domain_context.NewField("linkCountryViewCounter", linkId))
		return newLinkCountryViewCounter, nil
	}

	linkCountryViewCounter := linkCountryViewCounterOptional.Get()
	linkCountryViewCounter.Increment(totalViews, uniqueViews)
	u.logger.Info(ctx, "LinkCountryViewCounterIncrementer - Run - Incremented link country view counter", shared_domain_context.NewField("linkCountryViewCounter", linkId), shared_domain_context.NewField("uniqueViews", linkCountryViewCounter.UniqueViews), shared_domain_context.NewField("totalViews", linkCountryViewCounter.TotalViews))

	err = u.repo.Update(ctx, *linkCountryViewCounter)
	if err != nil {
		u.logger.Error(ctx, "LinkCountryViewCounterIncrementer - Run - Error updating link country view counter", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	u.logger.Info(ctx, "LinkCountryViewCounterIncrementer - Run - Successfully updated link country view counter", shared_domain_context.NewField("linkCountryViewCounter", linkId))

	return linkCountryViewCounter, nil
}
