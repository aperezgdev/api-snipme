package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkAnalyticsRemover_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid link id", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		remover := NewLinkAnalyticsRemover(logger, repo)

		idLink := "00000000-0000-0000-0000-000000000000"
		domainIdLink, _ := shared_domain_context.ParseID(idLink)

		repo.On("RemoveByLink", mock.Anything, domainIdLink).Return(nil)

		err := remover.Run(context.Background(), idLink)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on repository remove error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		remover := NewLinkAnalyticsRemover(logger, repo)

		idLink := "00000000-0000-0000-0000-000000000000"
		domainIdLink, _ := shared_domain_context.ParseID(idLink)

		repo.On("RemoveByLink", mock.Anything, domainIdLink).Return(errors.New("database error"))

		err := remover.Run(context.Background(), idLink)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid link id", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkAnalyticsRepositoryMock{}
		remover := NewLinkAnalyticsRemover(logger, repo)

		invalidId := "not-a-valid-id"

		err := remover.Run(context.Background(), invalidId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
