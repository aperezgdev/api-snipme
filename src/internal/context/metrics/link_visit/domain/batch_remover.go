package domain

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkVisitBatchRemover struct {
	logger shared_domain_context.Logger
	repo   LinkVisitRepository
}

func NewLinkVisitBatchRemover(logger shared_domain_context.Logger, repo LinkVisitRepository) *LinkVisitBatchRemover {
	return &LinkVisitBatchRemover{
		logger: logger,
		repo:   repo,
	}
}

func (r *LinkVisitBatchRemover) Run(ctx context.Context, linkVisitsIds []shared_domain_context.Id) error {
	r.logger.Info(ctx, "LinkVisitBatchRemover - Run - Removing link visits in batch", shared_domain_context.NewField("count", len(linkVisitsIds)))
	err := r.repo.RemoveAll(ctx, linkVisitsIds)
	if err != nil {
		r.logger.Error(ctx, "LinkVisitBatchRemover - Run - Error removing link visits in batch", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "LinkVisitBatchRemover - Run - Link visits removed successfully in batch", shared_domain_context.NewField("count", len(linkVisitsIds)))
	return nil
}
