package application

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/stretchr/testify/mock"
)

func TestShortLinkCreator_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid short link", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		creator := NewShortLinkCreator(logger, repo, eventBus)

		originalLink := "https://example.com/some/long/path"
		clientId := "00000000-0000-0000-0000-000000000000"

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.ShortLink")).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		shortLink, err := creator.Run(context.Background(), originalLink, clientId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(shortLink.OriginalRoute) != originalLink {
			t.Errorf("Expected original link %s, got %s", originalLink, string(shortLink.OriginalRoute))
		}

		repo.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})

	t.Run("Run fails on repository save error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		creator := NewShortLinkCreator(logger, repo, eventBus)

		originalLink := "https://example.com/some/long/path"
		clientId := "00000000-0000-0000-0000-000000000000"

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.ShortLink")).Return(errors.New("database error"))

		_, err := creator.Run(context.Background(), originalLink, clientId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid original link", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		creator := NewShortLinkCreator(logger, repo, eventBus)

		originalLink := "invalid-url"
		clientId := "00000000-0000-0000-0000-000000000000"

		_, err := creator.Run(context.Background(), originalLink, clientId)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Run success with empty client ID", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		creator := NewShortLinkCreator(logger, repo, eventBus)

		originalLink := "https://example.com/some/long/path"
		clientId := ""

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.ShortLink")).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		shortLink, err := creator.Run(context.Background(), originalLink, clientId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(shortLink.OriginalRoute) != originalLink {
			t.Errorf("Expected original link %s, got %s", originalLink, string(shortLink.OriginalRoute))
		}

		repo.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})
}
