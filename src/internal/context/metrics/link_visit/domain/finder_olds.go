package domain

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkVisitOldsFinder struct {
	logger shared_domain_context.Logger
	repo   LinkVisitRepository
}

func NewLinkVisitOldsFinder(logger shared_domain_context.Logger, repo LinkVisitRepository) *LinkVisitOldsFinder {
	return &LinkVisitOldsFinder{
		logger: logger,
		repo:   repo,
	}
}

func (f *LinkVisitOldsFinder) Run(ctx context.Context) ([]LinkVisit, error) {
	f.logger.Info(ctx, "LinkVisitOldsFinder - Run - Finding old link visits")
	linkVisits, err := f.repo.FindOlds(ctx)
	if err != nil {
		f.logger.Error(ctx, "LinkVisitOldsFinder - Run - Error finding old link visits", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}
	f.logger.Info(ctx, "LinkVisitOldsFinder - Run - Found old link visits", shared_domain_context.NewField("count", len(linkVisits)))
	return linkVisits, nil
}
