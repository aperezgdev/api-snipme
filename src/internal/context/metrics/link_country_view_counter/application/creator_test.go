package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkCountryViewCounterCreator_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid linkID and countryCode", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		creator := NewLinkCountryViewCounterCreator(logger, repo)

		linkID := "00000000-0000-0000-0000-000000000001"
		countryCode := "US"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(nil)

		counter, err := creator.Run(context.Background(), linkID, countryCode)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if counter.LinkId.String() != linkID {
			t.Errorf("Expected linkID %s, got %s", linkID, counter.LinkId.String())
		}
		if string(counter.CountryCode) != countryCode {
			t.Errorf("Expected countryCode %s, got %s", countryCode, counter.CountryCode)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid linkID or countryCode", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		creator := NewLinkCountryViewCounterCreator(logger, repo)

		invalidLinkID := ""
		countryCode := "US"

		_, err := creator.Run(context.Background(), invalidLinkID, countryCode)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Run fails on repository save error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		creator := NewLinkCountryViewCounterCreator(logger, repo)

		linkID := "00000000-0000-0000-0000-000000000001"
		countryCode := "US"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(errors.New("database error"))

		_, err := creator.Run(context.Background(), linkID, countryCode)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
