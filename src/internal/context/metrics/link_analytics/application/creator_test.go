package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkAnalyticsCreator_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid linkId", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		creator := NewLinkAnalyticsCreator(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(nil)

		result, err := creator.Run(context.Background(), linkId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result.Id.String() == "" {
			t.Errorf("Expected non-empty Id, got empty")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid linkId", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		creator := NewLinkAnalyticsCreator(logger, repo)

		invalidLinkId := ""

		_, err := creator.Run(context.Background(), invalidLinkId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Run fails on repository save error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		creator := NewLinkAnalyticsCreator(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkAnalytics")).Return(errors.New("database error"))

		_, err := creator.Run(context.Background(), linkId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
