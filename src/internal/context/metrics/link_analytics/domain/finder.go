package domain

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkAnalyticsFinder struct {
	logger shared_domain_context.Logger
	repo   LinkAnalyticsRepository
}

func NewLinkAnalyticsFinder(
	logger shared_domain_context.Logger,
	repo LinkAnalyticsRepository,
) *LinkAnalyticsFinder {
	return &LinkAnalyticsFinder{
		logger: logger,
		repo:   repo,
	}
}

func (f *LinkAnalyticsFinder) Run(ctx context.Context, linkId string) (*LinkAnalytics, error) {
	f.logger.Info(ctx, "LinkAnalyticsFinder - Run - Finding LinkAnalytics", shared_domain_context.NewField("linkId", linkId))

	idVo, err := shared_domain_context.ParseID(linkId)
	if err != nil {
		f.logger.Error(ctx, "LinkAnalyticsFinder - Run - Error parsing linkId", shared_domain_context.NewField("error", err.Error()))
		return &LinkAnalytics{}, shared_domain_context.NewValidationError("linkId", "invalid linkId format")
	}

	linkAnalytics, err := f.repo.FindByLinkId(ctx, idVo)
	if err != nil {
		f.logger.Error(ctx, "LinkAnalyticsFinder - Run - Error finding LinkAnalytics", shared_domain_context.NewField("error", err.Error()))
		return &LinkAnalytics{}, err
	}

	if !linkAnalytics.IsPresent() {
		f.logger.Info(ctx, "LinkAnalyticsFinder - Run - LinkAnalytics not found", shared_domain_context.NewField("linkId", linkId))
		return &LinkAnalytics{}, shared_domain_context.NewNotFoundError("LinkAnalytics not found for linkId: " + linkId)
	}

	f.logger.Info(ctx, "LinkAnalyticsFinder - Run - LinkAnalytics found successfully", shared_domain_context.NewField("linkId", linkId))
	return linkAnalytics.Get(), nil
}
