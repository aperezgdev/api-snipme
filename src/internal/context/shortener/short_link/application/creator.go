package application

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
)

type ShortLinkCreator struct {
	looger   shared_domain.Logger
	repo     domain.ShortLinkRepository
	eventBus shared_domain.EventBus
}

func NewShortLinkCreator(logger shared_domain.Logger, repo domain.ShortLinkRepository, eventBus shared_domain.EventBus) *ShortLinkCreator {
	return &ShortLinkCreator{
		looger:   logger,
		repo:     repo,
		eventBus: eventBus,
	}
}

func (c ShortLinkCreator) Run(ctx context.Context, orignalLink, clientId string) (*domain.ShortLink, error) {
	c.looger.Info(ctx, "ShortLinkCreator - Run - Params into: ", shared_domain.NewField("originalLink", orignalLink), shared_domain.NewField("clientId", clientId))
	shortLink, err := domain.NewShortLink(orignalLink, clientId)
	if err != nil {
		c.looger.Error(ctx, "ShortLinkCreator - Run - Error creating short link", shared_domain.NewField("error", err.Error()))
		return nil, err
	}

	err = c.repo.Save(ctx, shortLink)
	if err != nil {
		c.looger.Error(ctx, "ShortLinkCreator - Run - Error saving short link", shared_domain.NewField("error", err.Error()))
		return nil, err
	}

	c.eventBus.Publish(ctx, shortLink.PullDomainEvents()...)

	return shortLink, nil
}
