package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkAnalyticsFinder struct {
	logger shared_domain_context.Logger
	finder domain.LinkAnalyticsFinder
}

func NewLinkAnalyticsFinder(
	logger shared_domain_context.Logger,
	repo domain.LinkAnalyticsRepository,
) *LinkAnalyticsFinder {
	return &LinkAnalyticsFinder{
		logger: logger,
		finder: *domain.NewLinkAnalyticsFinder(logger, repo),
	}
}

func (f *LinkAnalyticsFinder) Run(ctx context.Context, idLink string) (*domain.LinkAnalytics, error) {
	f.logger.Info(ctx, "LinkAnalyticsFinder - Run - Params into: ", shared_domain_context.NewField("idLink", idLink))
	linkAnalytics, err := f.finder.Run(ctx, idLink)
	if err != nil {
		f.logger.Error(ctx, "LinkAnalyticsFinder - Run - Error finding link analytics", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}
	f.logger.Info(ctx, "LinkAnalyticsFinder - Run - Successfully found link analytics for link", shared_domain_context.NewField("idLink", idLink))
	return linkAnalytics, nil
}
