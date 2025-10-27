package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	short_link_domain "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/stretchr/testify/mock"
)

func TestCreatorOnShortLinkCreated_On(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("On success when LinkAnalytics is created", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		creator := NewCreatorOnShortLinkCreated(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000001"
		event := short_link_domain.NewShortLinkCreated(linkId)

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(nil)

		err := creator.On(context.Background(), event)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("On fails when LinkAnalytics creation fails", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		creator := NewCreatorOnShortLinkCreated(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000002"
		event := short_link_domain.NewShortLinkCreated(linkId)

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(errors.New("database error"))

		err := creator.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
