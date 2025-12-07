package domain

import (
	"context"
	"errors"
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestLinkCountryViewCounterIncrementer_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run counter not found should create a new one", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000001"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := shared_domain.NewCountryCode(countryCode)
		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(pkg.EmptyOptional[*LinkCountryViewCounter](), nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(nil)
		
		result, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)

		if err != nil {
			t.Fatalf("Expected error, got nil")
		}
		if result == nil {
			t.Fatalf("Expected result, got nil")
		}

		if uint(result.UniqueViews) != uniqueViews || uint(result.TotalViews) != totalViews {
			t.Errorf("Expected views (unique: %d, total: %d), got (unique: %d, total: %d)", uniqueViews, totalViews, result.UniqueViews, result.TotalViews)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when finder returns error", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000002"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := shared_domain.NewCountryCode(countryCode)
		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(pkg.EmptyOptional[*LinkCountryViewCounter](), errors.New("db error"))

		_, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run success when counter found and updated", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000003"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := shared_domain.NewCountryCode(countryCode)
		counter := LinkCountryViewCounter{
			LinkId:      domainId,
			CountryCode: countryCodeVO,
			UniqueViews: 5,
			TotalViews:  10,
		}
		opt := pkg.Some(&counter)

		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(opt, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(nil)

		result, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if uint(result.UniqueViews) != 15 || uint(result.TotalViews) != 30 {
			t.Errorf("Expected updated views (unique: %d, total: %d), got (unique: %d, total: %d)", 15, 30, result.UniqueViews, result.TotalViews)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails when repo update returns error", func(t *testing.T) {
		t.Parallel()
		repo := &LinkCountryViewCounterRepositoryMock{}
		finder := NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000004"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := shared_domain.NewCountryCode(countryCode)
		counter := LinkCountryViewCounter{
			LinkId:      domainId,
			CountryCode: countryCodeVO,
			UniqueViews: 5,
			TotalViews:  10,
		}
		opt := pkg.Some(&counter)

		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(opt, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(errors.New("update error"))

		_, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})
}
