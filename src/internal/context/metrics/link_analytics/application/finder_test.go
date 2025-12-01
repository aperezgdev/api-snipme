package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	link_country_view_counter_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestLinkAnalyticsFinder(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run returns link analytics when found", func(t *testing.T) {
		repo := &domain.LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainLinkId, _ := shared_domain_context.ParseID(linkId)

		analyticsId, _ := shared_domain_context.NewID()
		linkAnalytics := domain.LinkAnalytics{
			Id:                        analyticsId,
			LinkId:                    domainLinkId,
			TotalViews:                shared_domain.NewLinkViewsCounter(10),
			UniqueViews:               shared_domain.NewLinkViewsCounter(5),
			LinkCountriesViewCounters: []link_country_view_counter_domain.LinkCountryViewCounter{},
			UpdateOn:                  shared_domain_context.NewUpdatedOn(),
		}

		repo.On("FindByLinkId", mock.Anything, domainLinkId).Return(pkg.Some(&linkAnalytics), nil)

		foundAnalytics, err := finder.Run(context.Background(), linkId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if foundAnalytics.LinkId != linkAnalytics.LinkId {
			t.Fatalf("Expected link ID %v, got %v", linkAnalytics.LinkId, foundAnalytics.LinkId)
		}
		if foundAnalytics.TotalViews != linkAnalytics.TotalViews {
			t.Fatalf("Expected total views %v, got %v", linkAnalytics.TotalViews, foundAnalytics.TotalViews)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run returns error when link analytics not found", func(t *testing.T) {
		repo := &domain.LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainLinkId, _ := shared_domain_context.ParseID(linkId)

		expectedError := errors.New("link analytics not found")
		repo.On("FindByLinkId", mock.Anything, domainLinkId).Return(pkg.EmptyOptional[*domain.LinkAnalytics](), expectedError)

		_, err := finder.Run(context.Background(), linkId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected error %v, got %v", expectedError, err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run returns error when linkId is invalid", func(t *testing.T) {
		repo := &domain.LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		invalidLinkId := "invalid-uuid"

		_, err := finder.Run(context.Background(), invalidLinkId)
		if err == nil {
			t.Fatalf("Expected error for invalid UUID, got nil")
		}

		repo.AssertNotCalled(t, "FindByLinkId")
	})
}