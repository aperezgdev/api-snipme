package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestLinkCountryViewCounterIncrementer_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run fails when counter not found", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		finder := domain.NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo, eventBus)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000001"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := domain.NewCountryCode(countryCode)
		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(pkg.EmptyOptional[domain.LinkCountryViewCounter](), nil)

		_, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run fails when finder returns error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		finder := domain.NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo, eventBus)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000002"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := domain.NewCountryCode(countryCode)
		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(pkg.EmptyOptional[domain.LinkCountryViewCounter](), errors.New("db error"))

		_, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run success when counter found and updated", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		finder := domain.NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo, eventBus)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000003"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := domain.NewCountryCode(countryCode)
		counter := domain.LinkCountryViewCounter{
			LinkId:      domainId,
			CountryCode: countryCodeVO,
			UniqueViews: 5,
			TotalViews:  10,
		}
		opt := pkg.Some(counter)

		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(opt, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		result, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if uint(result.UniqueViews) != uint(counter.UniqueViews)+uniqueViews || uint(result.TotalViews) != uint(counter.TotalViews)+totalViews {
			t.Errorf("Expected updated views (unique: %d, total: %d), got (unique: %d, total: %d)", uint(counter.UniqueViews)+uniqueViews, uint(counter.TotalViews)+totalViews, result.UniqueViews, result.TotalViews)
		}

		repo.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})

	t.Run("Run fails when repo update returns error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		finder := domain.NewLinkCountryViewCounterFinder(logger, repo)
		incrementer := NewLinkCountryViewCounterIncrementer(logger, repo, eventBus)
		incrementer.finder = finder

		linkId := "00000000-0000-0000-0000-000000000004"
		countryCode := "US"
		uniqueViews := uint(10)
		totalViews := uint(20)

		domainId, _ := shared_domain_context.ParseID(linkId)
		countryCodeVO, _ := domain.NewCountryCode(countryCode)
		counter := domain.LinkCountryViewCounter{
			LinkId:      domainId,
			CountryCode: countryCodeVO,
			UniqueViews: 5,
			TotalViews:  10,
		}
		opt := pkg.Some(counter)

		repo.On("Find", mock.Anything, domainId, countryCodeVO).Return(opt, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(errors.New("update error"))

		_, err := incrementer.Run(context.Background(), linkId, countryCode, uniqueViews, totalViews)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})
}
