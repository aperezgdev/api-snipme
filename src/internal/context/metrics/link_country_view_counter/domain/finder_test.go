package domain

import (
	"context"
	"errors"
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_context_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestLinkCountryViewCounter(t *testing.T) {
	logger := shared_context_domain.DummyLogger{}

	t.Run("Run successfully should find link country view counter", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)

		idLink := "00000000-0000-0000-0000-000000000000"
		countryCode := "US"

		expectedLinkCountryViewCounter := &LinkCountryViewCounter{
			CountryCode: shared_domain.CountryCode(countryCode),
			TotalViews:       10,
		}

		repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(pkg.Some(expectedLinkCountryViewCounter), nil)

		result, err := finder.Run(context.Background(), idLink, countryCode)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !result.IsPresent() || result.Get().TotalViews != expectedLinkCountryViewCounter.TotalViews {
			t.Errorf("Expected link country view counter with views %d, got %v", expectedLinkCountryViewCounter.TotalViews, result)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when idLink is invalid", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)

		invalidIdLink := "invalid-id"
		countryCode := "US"

		_, err := finder.Run(context.Background(), invalidIdLink, countryCode)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when country code is invalid", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)

		idLink := "507f1f77bcf86cd799439011"
		invalidCountryCode := "INVALID"

		_, err := finder.Run(context.Background(), idLink, invalidCountryCode)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when repository returns an error", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)

		idLink := "00000000-0000-0000-0000-000000000000"
		countryCode := "US"

		repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*LinkCountryViewCounter](), errors.New("database error"))

		_, err := finder.Run(context.Background(), idLink, countryCode)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}