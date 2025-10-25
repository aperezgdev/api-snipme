package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkVisitCreator struct {
	logger shared_domain_context.Logger
	repo   domain.LinkVisitRepository
}

func NewLinkVisitCreator(
	logger shared_domain_context.Logger,
	repo domain.LinkVisitRepository,
) *LinkVisitCreator {
	return &LinkVisitCreator{
		logger: logger,
		repo:   repo,
	}
}

func (lc LinkVisitCreator) Run(ctx context.Context, linkId, ip, userAgent string) (*domain.LinkVisit, error) {
	lc.logger.Info(ctx, "LinkVisitCreator - Run: Creating LinkVisit", shared_domain_context.NewField("linkId", linkId), shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("userAgent", userAgent))
	linkVisit, err := domain.NewLinkVisit(linkId, ip, userAgent)
	if err != nil {
		lc.logger.Error(ctx, "LinkVisitCreator - Run: Error creating LinkVisit entity", shared_domain_context.NewField("error", err), shared_domain_context.NewField("linkId", linkId), shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("userAgent", userAgent))
		return nil, err
	}

	if err := lc.repo.Save(ctx, *linkVisit); err != nil {
		lc.logger.Error(ctx, "LinkVisitCreator - Run: Error saving LinkVisit entity", shared_domain_context.NewField("error", err), shared_domain_context.NewField("linkVisit", linkVisit))
		return nil, err
	}

	lc.logger.Info(ctx, "LinkVisitCreator - Run: LinkVisit created successfully", shared_domain_context.NewField("linkVisitId", linkVisit.Id.String()))
	return linkVisit, nil
}
