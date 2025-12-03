package application

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
)

type PublicShortLinkCreator struct {
	logger shared_domain_context.Logger
	repo  domain.ShortLinkRepository
	eventBus shared_domain_context.EventBus
}

func NewPublicShortLinkCreator(logger shared_domain_context.Logger, repo domain.ShortLinkRepository, eventBus shared_domain_context.EventBus) *PublicShortLinkCreator {
	return &PublicShortLinkCreator{
		logger: logger,
		repo:  repo,
		eventBus: eventBus,
	}
}

func (c PublicShortLinkCreator) Run(ctx context.Context, originalLink string) (*domain.ShortLink, error) {
	c.logger.Info(ctx, "PublicShortLinkCreator - Run - Params into: ", shared_domain_context.NewField("originalLink", originalLink))
	shortLink, err := domain.NewPublicShortLink(originalLink)
	if err != nil {
		c.logger.Error(ctx, "PublicShortLinkCreator - Run - Error creating short link", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	err = c.repo.Save(ctx, shortLink)
	if err != nil {
		c.logger.Error(ctx, "PublicShortLinkCreator - Run - Error saving short link", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	c.eventBus.Publish(ctx, shortLink.PullDomainEvents()...)

	return shortLink, nil
}