package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	domain_link_visit "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type UpdaterOnLinkVisitProcessed struct {
	logger shared_domain_context.Logger
	finder domain.LinkAnalyticsFinder
	repo   domain.LinkAnalyticsRepository
}

func NewUpdaterOnLinkVisitProcessed(
	logger shared_domain_context.Logger,
	repo domain.LinkAnalyticsRepository,
) *UpdaterOnLinkVisitProcessed {
	return &UpdaterOnLinkVisitProcessed{
		logger: logger,
		finder: *domain.NewLinkAnalyticsFinder(logger, repo),
		repo:   repo,
	}
}

func (u *UpdaterOnLinkVisitProcessed) On(
	ctx context.Context,
	event shared_domain_context.DomainEvent,
) error {
	eventData, ok := event.(domain_link_visit.LinkVisitsProcessed)
	if !ok {
		u.logger.Error(ctx, "UpdaterOnLinkVisitProcessed - On - Invalid event type", shared_domain_context.NewField("expected", domain_link_visit.LinkVisitsProcessedEventName), shared_domain_context.NewField("received", event.Name()))
		return nil
	}

	u.logger.Info(ctx, "UpdaterOnLinkVisitProcessed - On - Params in",
		shared_domain_context.NewField("linkId", eventData.AggregateId),
		shared_domain_context.NewField("visitsCount", eventData.TotalViews),
		shared_domain_context.NewField("uniqueViews", eventData.UniqueViews),
	)

	linkAnalytics, err := u.finder.Run(ctx, eventData.AggregateId)
	if err != nil {
		u.logger.Error(ctx, "UpdaterOnLinkVisitProcessed - On - Error finding LinkAnalytics", shared_domain_context.NewField("error", err.Error()))
		return err
	}

	linkAnalytics.IncrementTotalViews(eventData.TotalViews)
	linkAnalytics.IncrementUniqueViews(eventData.UniqueViews)

	err = u.repo.Update(ctx, *linkAnalytics)
	if err != nil {
		u.logger.Error(ctx, "UpdaterOnLinkVisitProcessed - On - Error saving LinkAnalytics", shared_domain_context.NewField("error", err.Error()))
		return err
	}

	u.logger.Info(ctx, "UpdaterOnLinkVisitProcessed - On - LinkAnalytics updated successfully", shared_domain_context.NewField("linkId", eventData.AggregateId))
	return nil
}
