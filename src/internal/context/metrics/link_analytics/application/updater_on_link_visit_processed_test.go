package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	link_visit_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	short_link_domain "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestUpdaterOnLinkVisitProcessed_On(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("On success when LinkAnalytics is updated with new views", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		updater := NewUpdaterOnLinkVisitProcessed(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000001"
		domainId, _ := shared_domain_context.ParseID(linkId)

		existingLinkAnalytics := &domain.LinkAnalytics{
			Id:          domainId,
			LinkId:      domainId,
			TotalViews:  shared_domain.NewLinkViewsCounter(10),
			UniqueViews: shared_domain.NewLinkViewsCounter(5),
		}

		event := link_visit_domain.NewLinkVisitsProcessedDomainEvent(linkId, 3, 2)

		repo.On("FindByLinkId", mock.Anything, domainId).Return(pkg.Some(existingLinkAnalytics), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(nil)

		err := updater.On(context.Background(), event)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)

		repo.AssertCalled(t, "Update", mock.Anything, mock.MatchedBy(func(la domain.LinkAnalytics) bool {
			return uint(la.TotalViews) == 13 && uint(la.UniqueViews) == 7
		}))
	})

	t.Run("On fails when LinkAnalytics is not found", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		updater := NewUpdaterOnLinkVisitProcessed(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000002"
		domainId, _ := shared_domain_context.ParseID(linkId)

		event := link_visit_domain.NewLinkVisitsProcessedDomainEvent(linkId, 3, 2)

		repo.On("FindByLinkId", mock.Anything, domainId).Return(pkg.EmptyOptional[*domain.LinkAnalytics](), errors.New("not found"))

		err := updater.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
		repo.AssertNotCalled(t, "Save")
	})

	t.Run("On fails when Save returns error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		updater := NewUpdaterOnLinkVisitProcessed(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000003"
		domainId, _ := shared_domain_context.ParseID(linkId)

		existingLinkAnalytics := &domain.LinkAnalytics{
			Id:          domainId,
			LinkId:      domainId,
			TotalViews:  shared_domain.NewLinkViewsCounter(10),
			UniqueViews: shared_domain.NewLinkViewsCounter(5),
		}

		event := link_visit_domain.NewLinkVisitsProcessedDomainEvent(linkId, 3, 2)

		repo.On("FindByLinkId", mock.Anything, domainId).Return(pkg.Some(existingLinkAnalytics), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(errors.New("database error"))

		err := updater.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("On returns nil when event type is invalid", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		updater := NewUpdaterOnLinkVisitProcessed(logger, repo)

		invalidEvent := short_link_domain.NewShortLinkCreated("some-id")

		err := updater.On(context.Background(), invalidEvent)
		if err != nil {
			t.Fatalf("Expected no error for invalid event type, got %v", err)
		}

		repo.AssertNotCalled(t, "FindByLinkId")
		repo.AssertNotCalled(t, "Update")
	})
}
