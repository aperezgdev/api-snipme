package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkAnalyticsCreator struct {
	logger shared_domain_context.Logger
	repo   domain.LinkAnalyticsRepository
}

func NewLinkAnalyticsCreator(logger shared_domain_context.Logger, repo domain.LinkAnalyticsRepository) *LinkAnalyticsCreator {
	return &LinkAnalyticsCreator{
		logger: logger,
		repo:   repo,
	}
}

func (c LinkAnalyticsCreator) Run(ctx context.Context, linkId string) (*domain.LinkAnalytics, error) {
	c.logger.Info(ctx, "LinkAnalyticsCreator - Run - Params into: ", shared_domain_context.NewField("linkId", linkId))
	linkAnalytics, err := domain.NewLinkAnalytics(linkId)
	if err != nil {
		c.logger.Error(ctx, "LinkAnalyticsCreator - Run - Error creating link analytics", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	err = c.repo.Save(ctx, *linkAnalytics)
	if err != nil {
		c.logger.Error(ctx, "LinkAnalyticsCreator - Run - Error saving link analytics", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}
	c.logger.Info(ctx, "LinkAnalyticsCreator - Run - Link analytics created successfully", shared_domain_context.NewField("linkAnalyticsId", linkAnalytics.Id.String()))

	return linkAnalytics, nil
}
