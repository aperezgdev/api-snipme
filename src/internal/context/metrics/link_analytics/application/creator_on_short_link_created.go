package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type CreatorOnShortLinkCreated struct {
	logger  shared_domain_context.Logger
	creator LinkAnalyticsCreator
}

func NewCreatorOnShortLinkCreated(
	logger shared_domain_context.Logger,
	repo domain.LinkAnalyticsRepository,
) *CreatorOnShortLinkCreated {
	return &CreatorOnShortLinkCreated{
		logger:  logger,
		creator: *NewLinkAnalyticsCreator(logger, repo),
	}
}

func (c CreatorOnShortLinkCreated) On(ctx context.Context, event shared_domain_context.DomainEvent) error {
	c.logger.Info(ctx, "CreatorOnShortLinkCreated - On - Received ShortLinkCreated event", shared_domain_context.NewField("aggregateID", event.AggregateID()))
	_, err := c.creator.Run(ctx, event.AggregateID())
	if err != nil {
		c.logger.Error(ctx, "CreatorOnShortLinkCreated - On - Error creating link analytics", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	return err
}
