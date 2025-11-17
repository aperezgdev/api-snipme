package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type ShortLinkFinder struct {
	logger shared_domain.Logger
	repo   ShortLinkRepository
}

func NewShortLinkFinder(logger shared_domain.Logger, repo ShortLinkRepository) *ShortLinkFinder {
	return &ShortLinkFinder{
		logger: logger,
		repo:   repo,
	}
}

func (f ShortLinkFinder) Run(ctx context.Context, id string) (*ShortLink, error) {
	f.logger.Info(ctx, "ShortLinkFinder - Run - Params into: ", shared_domain.NewField("id", id))
	domainId, err := shared_domain.ParseID(id)
	if err != nil {
		f.logger.Error(ctx, "ShortLinkFinder - Run - Error parsing id", shared_domain.NewField("error", err.Error()))
		return nil, err
	}
	shortLink, err := f.repo.FindById(ctx, domainId)
	if err != nil {
		f.logger.Error(ctx, "ShortLinkFinder - Run - Error finding short link", shared_domain.NewField("error", err.Error()))
		return nil, err
	}
	if !shortLink.IsPresent() {
		f.logger.Info(ctx, "ShortLinkFinder - Run - No short link found", shared_domain.NewField("id", id))
		return nil, shared_domain.NewNotFoundError("Short link not found")
	}

	shortLinkValue := shortLink.Get()
	f.logger.Info(ctx, "ShortLinkFinder - Run - Short link found", shared_domain.NewField("shortLinkId", shortLinkValue.Id.String()))

	return shortLinkValue, nil
}
