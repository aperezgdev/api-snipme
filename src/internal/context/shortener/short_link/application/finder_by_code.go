package application

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
)

type ShortLinkFinderByCode struct {
	logger shared_domain_context.Logger
	repo   domain.ShortLinkRepository
}

func NewShortLinkFinderByCode(
	logger shared_domain_context.Logger,
	repo domain.ShortLinkRepository,
) *ShortLinkFinderByCode {
	return &ShortLinkFinderByCode{
		logger: logger,
		repo:   repo,
	}
}

func (sf ShortLinkFinderByCode) Run(ctx context.Context, code string) (*domain.ShortLink, error) {
	sf.logger.Info(ctx, "ShortLinkFinderByCode - Run: Finding short link by code", shared_domain_context.NewField("code", code))

	shortLinkOpt, err := sf.repo.FindByCode(ctx, domain.ShortLinkCode(code))
	if err != nil {
		sf.logger.Error(ctx, "ShortLinkFinderByCode - Run: Error finding short link by code", shared_domain_context.NewField("error", err), shared_domain_context.NewField("code", code))
		return nil, err
	}

	if !shortLinkOpt.IsPresent() {
		sf.logger.Info(ctx, "ShortLinkFinderByCode - Run: No short link found for the given code", shared_domain_context.NewField("code", code))
		return nil, shared_domain_context.NewNotFoundError("Short link not found")
	}

	shortLink := shortLinkOpt.Get()

	sf.logger.Info(ctx, "ShortLinkFinderByCode - Run: Short link found", shared_domain_context.NewField("shortLinkId", shortLink.Id.String()), shared_domain_context.NewField("code", code))

	return shortLink, nil
}
