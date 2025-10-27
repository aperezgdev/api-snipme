package application

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/stretchr/testify/mock"
)

func TestShortLinkRemover_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid id", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		remover := NewShortLinkRemover(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		repo.On("Remove", mock.Anything, domainId).Return(nil)

		err := remover.Run(context.Background(), id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on repository remove error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		remover := NewShortLinkRemover(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		repo.On("Remove", mock.Anything, domainId).Return(errors.New("database error"))

		err := remover.Run(context.Background(), id)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid id", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		remover := NewShortLinkRemover(logger, repo)

		invalidId := "not-a-valid-id"

		err := remover.Run(context.Background(), invalidId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
