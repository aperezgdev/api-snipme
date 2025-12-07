package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestCountryFinder(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success when country is found by IP", func(t *testing.T) {
		t.Parallel()
		repo := &domain.CountryRepositoryMock{}
		finder := NewCountryFinder(logger, repo)

		ip := "1.1.1.1"
		expectedCountry := &domain.Country{
			ISOCode: shared_domain.CountryCode("US"),
		}

		repo.On("FindByIp", mock.Anything, mock.Anything).Return(pkg.Some(expectedCountry), nil)

		country, err := finder.Run(context.Background(), ip)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if country.ISOCode != expectedCountry.ISOCode {
			t.Errorf("Expected country ISO code %s, got %s", expectedCountry.ISOCode, country.ISOCode)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when country is not found by IP", func(t *testing.T) {
		t.Parallel()
		repo := &domain.CountryRepositoryMock{}
		finder := NewCountryFinder(logger, repo)

		ip := "1.1.1.1"

		repo.On("FindByIp", mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.Country](), nil)

		_, err := finder.Run(context.Background(), ip)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		var notFoundErr shared_domain_context.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Fatalf("Expected NotFoundError, got %v", err)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when repository returns an error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.CountryRepositoryMock{}
		finder := NewCountryFinder(logger, repo)

		ip := "1.1.1.1"

		repo.On("FindByIp", mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.Country](), errors.New("database error"))

		_, err := finder.Run(context.Background(), ip)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when ip is not valid", func(t *testing.T) {
		t.Parallel()
		repo := &domain.CountryRepositoryMock{}
		finder := NewCountryFinder(logger, repo)

		invalidIp := "invalid-ip"

		_, err := finder.Run(context.Background(), invalidIp)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if !errors.As(err, &shared_domain_context.ValidationError{}) {
			t.Fatalf("Expected ValidationError, got %v", err)
		}

		repo.AssertNotCalled(t, "FindByIp")
	})
}
