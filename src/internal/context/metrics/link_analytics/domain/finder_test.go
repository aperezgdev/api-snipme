package domain

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestLinkAnalyticsFinder(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}
	t.Run("Run returns link analytics when found", func(t *testing.T) {
		repo := &LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(linkId)

		linkAnalytics := &LinkAnalytics{
				Id:     domainId,
				LinkId: domainId,
		}

		repo.On("FindByLinkId", mock.Anything, domainId).Return(pkg.Some(linkAnalytics), nil)

		foundLinkAnalytics, err := finder.Run(context.Background(), linkId)
		if err != nil {
				t.Fatalf("Expected no error, got %v", err)
		}
		if foundLinkAnalytics.LinkId != linkAnalytics.LinkId {
				t.Fatalf("Expected link ID %v, got %v", linkAnalytics.LinkId, foundLinkAnalytics.LinkId)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run returns error when not found", func(t *testing.T) {
		repo := &LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(linkId)

		repo.On("FindByLinkId", mock.Anything, domainId).Return(pkg.EmptyOptional[*LinkAnalytics](), errors.New("not found"))

		_, err := finder.Run(context.Background(), linkId)
		if err == nil {
				t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run returns error when linkId is invalid", func(t *testing.T) {
		repo := &LinkAnalyticsRepositoryMock{}
		finder := NewLinkAnalyticsFinder(logger, repo)

		invalidLinkId := "invalid-uuid"

		_, err := finder.Run(context.Background(), invalidLinkId)
		if err == nil {
				t.Fatalf("Expected error for invalid UUID, got nil")
		}

		repo.AssertNotCalled(t, "FindByLinkId")
	})
}