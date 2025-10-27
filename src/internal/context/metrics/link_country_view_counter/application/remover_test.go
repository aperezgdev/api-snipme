package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkViewsCounterRemover_Run(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("Run success on valid linkId", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		remover := NewLinkViewsCounterRemover(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain.ParseID(linkId)

		repo.On("RemoveByLink", mock.Anything, domainId).Return(nil)

		err := remover.Run(context.Background(), linkId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on repository remove error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		remover := NewLinkViewsCounterRemover(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain.ParseID(linkId)

		repo.On("RemoveByLink", mock.Anything, domainId).Return(errors.New("database error"))

		err := remover.Run(context.Background(), linkId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid linkId", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		remover := NewLinkViewsCounterRemover(logger, repo)

		invalidId := "not-a-valid-id"

		err := remover.Run(context.Background(), invalidId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
