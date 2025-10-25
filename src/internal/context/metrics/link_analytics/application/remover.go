package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkAnalyticsRemover struct {
	logger shared_domain_context.Logger
	repo   domain.LinkAnalyticsRepository
}

func NewLinkAnalyticsRemover(logger shared_domain_context.Logger, repo domain.LinkAnalyticsRepository) *LinkAnalyticsRemover {
	return &LinkAnalyticsRemover{
		logger: logger,
		repo:   repo,
	}
}

func (r LinkAnalyticsRemover) Run(ctx context.Context, idLink string) error {
	r.logger.Info(ctx, "LinkAnalyticsRemover - Run - Params into: ", shared_domain_context.NewField("idLink", idLink))
	domainIdLink, err := shared_domain_context.ParseID(idLink)
	if err != nil {
		r.logger.Error(ctx, "LinkAnalyticsRemover - Run - Error parsing idLink", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	err = r.repo.RemoveByLink(ctx, domainIdLink)
	if err != nil {
		r.logger.Error(ctx, "LinkAnalyticsRemover - Run - Error removing link analytics", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "LinkAnalyticsRemover - Run - Successfully removed link analytics for link", shared_domain_context.NewField("idLink", idLink))

	return nil
}
