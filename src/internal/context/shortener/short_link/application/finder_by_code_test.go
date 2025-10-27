package application

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestShortLinkFinderByCode_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success when short link is found", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		finder := NewShortLinkFinderByCode(logger, repo)

		code := "abc123"
		expectedShortLink := &domain.ShortLink{
			Code:          domain.ShortLinkCode(code),
			OriginalRoute: "https://example.com",
		}

		opt := pkg.Some(expectedShortLink)
		repo.On("FindByCode", mock.Anything, domain.ShortLinkCode(code)).Return(opt, nil)

		shortLink, err := finder.Run(context.Background(), code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if shortLink.Code != domain.ShortLinkCode(code) {
			t.Errorf("Expected code %s, got %s", code, shortLink.Code)
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run fails when repository returns error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		finder := NewShortLinkFinderByCode(logger, repo)

		code := "abc123"
		repo.On("FindByCode", mock.Anything, domain.ShortLinkCode(code)).Return(pkg.EmptyOptional[*domain.ShortLink](), errors.New("db error"))

		_, err := finder.Run(context.Background(), code)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run fails when short link not found", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ShortLinkRepositoryMock{}
		finder := NewShortLinkFinderByCode(logger, repo)

		code := "notfound"
		opt := pkg.EmptyOptional[*domain.ShortLink]() // Not present
		repo.On("FindByCode", mock.Anything, domain.ShortLinkCode(code)).Return(opt, nil)

		_, err := finder.Run(context.Background(), code)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})
}
