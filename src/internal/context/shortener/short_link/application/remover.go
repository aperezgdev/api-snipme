package application

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
)

type ShortLinkRemover struct {
	logger shared_domain.Logger
	repo   domain.ShortLinkRepository
}

func NewShortLinkRemover(logger shared_domain.Logger, repo domain.ShortLinkRepository) *ShortLinkRemover {
	return &ShortLinkRemover{
		logger: logger,
		repo:   repo,
	}
}

func (r ShortLinkRemover) Run(ctx context.Context, id string) error {
	r.logger.Info(ctx, "ShortLinkRemover - Run - Params into: ", shared_domain.NewField("id", id))
	domainId, err := shared_domain.ParseID(id)
	if err != nil {
		r.logger.Error(ctx, "ShortLinkRemover - Run - Error parsing id", shared_domain.NewField("error", err.Error()))
		return err
	}
	err = r.repo.Remove(ctx, domainId)
	if err != nil {
		r.logger.Error(ctx, "ShortLinkRemover - Run - Error removing short link", shared_domain.NewField("error", err.Error()))
		return err
	}
	return nil
}
