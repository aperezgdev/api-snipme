package domain

import (
	"context"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestShortLinkFinder(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("FindById returns short link when found", func(t *testing.T) {
		t.Parallel()
		repo := &ShortLinkRepositoryMock{}
		finder := NewShortLinkFinder(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		shortLink := ShortLink{
			Id: domainId,
		}

		repo.On("FindById", mock.Anything, domainId).Return(pkg.Some(&shortLink), nil)

		foundShortLink, err := finder.Run(context.Background(), id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if foundShortLink.Id != shortLink.Id {
			t.Fatalf("Expected short link ID %v, got %v", shortLink.Id, foundShortLink.Id)
		}

		repo.AssertExpectations(t)
	})

	t.Run("FindById returns error when not found", func(t *testing.T) {
		t.Parallel()
		repo := &ShortLinkRepositoryMock{}
		finder := NewShortLinkFinder(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		repo.On("FindById", mock.Anything, domainId).Return(pkg.EmptyOptional[*ShortLink](), nil)

		_, err := finder.Run(context.Background(), id)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("FindById returns error on invalid ID", func(t *testing.T) {
		t.Parallel()
		repo := &ShortLinkRepositoryMock{}
		finder := NewShortLinkFinder(logger, repo)

		id := "invalid-uuid"

		_, err := finder.Run(context.Background(), id)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("FindById returns error on repository error", func(t *testing.T) {
		t.Parallel()
		repo := &ShortLinkRepositoryMock{}
		finder := NewShortLinkFinder(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		repo.On("FindById", mock.Anything, domainId).Return(pkg.EmptyOptional[*ShortLink](), shared_domain_context.NewNotFoundError("database error"))

		_, err := finder.Run(context.Background(), id)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
