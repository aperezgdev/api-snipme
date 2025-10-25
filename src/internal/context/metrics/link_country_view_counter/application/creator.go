package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkCountryViewCounterCreator struct {
	logger shared_domain.Logger
	repo   domain.LinkCountryViewCounterRepository
}

func NewLinkCountryViewCounterCreator(logger shared_domain.Logger, repo domain.LinkCountryViewCounterRepository) *LinkCountryViewCounterCreator {
	return &LinkCountryViewCounterCreator{
		logger: logger,
		repo:   repo,
	}
}

func (c LinkCountryViewCounterCreator) Run(ctx context.Context, linkID, countryCode string) (*domain.LinkCountryViewCounter, error) {
	c.logger.Info(ctx, "LinkCountryViewCounterCreator - Run - Params into: ", shared_domain.NewField("linkID", linkID), shared_domain.NewField("countryCode", countryCode))
	counter, err := domain.NewLinkCountryViewCounter(linkID, countryCode)
	if err != nil {
		c.logger.Error(ctx, "LinkCountryViewCounterCreator - Run - Error creating link country view counter", shared_domain.NewField("error", err.Error()))
		return nil, err
	}

	err = c.repo.Save(ctx, *counter)
	if err != nil {
		c.logger.Error(ctx, "LinkCountryViewCounterCreator - Run - Error saving link country view counter", shared_domain.NewField("error", err.Error()))
		return nil, err
	}

	return counter, nil
}
