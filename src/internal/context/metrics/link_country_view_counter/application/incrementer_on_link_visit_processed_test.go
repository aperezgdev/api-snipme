package application

import (
	"context"
	"errors"
	"testing"

	geo_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	domain_link_visit "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestIncremeterOnLinkVisitProcessed_On(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("returns error if geoRepo fails", func(t *testing.T) {
		geoRepo := &geo_domain.CountryRepositoryMock{}
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		geoRepo.On("FindByIp", mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*geo_domain.Country](), errors.New("geo error"))
		handler := NewIncremeterOnLinkVisitProcessed(logger, geoRepo, repo)
		linkId := "019af069-4ef2-7068-9dbb-48f7e1b9dd37"
		event := domain_link_visit.NewLinkVisitsProcessedDomainEvent(
			linkId,
			1,
			1,
			[]domain_link_visit.IpsVisits{{Ip: "102.177.191.0", TotalViews: 1, UniqueViews: 1}},
		)
		err := handler.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error from geoRepo, got nil")
		}
	})

	t.Run("returns error if repo.Save fails", func(t *testing.T) {
		geoRepo := &geo_domain.CountryRepositoryMock{}
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		geoRepo.On("FindByIp", mock.Anything, mock.Anything).Return(pkg.Some(&geo_domain.Country{ISOCode: "ES"}), nil)
		repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.LinkCountryViewCounter](), nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(errors.New("repo error"))
		handler := NewIncremeterOnLinkVisitProcessed(logger, geoRepo, repo)
		linkId := "019af069-4ef2-7068-9dbb-48f7e1b9dd37"
		event := domain_link_visit.NewLinkVisitsProcessedDomainEvent(
			linkId,
			1,
			1,
			[]domain_link_visit.IpsVisits{{Ip: "102.177.191.0", TotalViews: 1, UniqueViews: 1}},
		)
		err := handler.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error from repo.Save, got nil")
		}
	})

	t.Run("increments counters for each country", func(t *testing.T) {
		geoRepo := &geo_domain.CountryRepositoryMock{}
		repo := &domain.LinkCountryViewCounterRepositoryMock{}
		geoRepo.On("FindByIp", mock.Anything, mock.MatchedBy(func(ip shared_domain.Ip) bool {
			return ip.String() == "102.177.191.0"
		})).Return(pkg.Some(&geo_domain.Country{ISOCode: "ES"}), nil)
		geoRepo.On("FindByIp", mock.Anything, mock.MatchedBy(func(ip shared_domain.Ip) bool {
			return ip.String() == "102.212.56.0"
		})).Return(pkg.Some(&geo_domain.Country{ISOCode: "ES"}), nil)
		geoRepo.On("FindByIp", mock.Anything, mock.MatchedBy(func(ip shared_domain.Ip) bool {
			return ip.String() == "102.129.65.0"
		})).Return(pkg.Some(&geo_domain.Country{ISOCode: "FR"}), nil)
		geoRepo.On("FindByIp", mock.Anything, mock.MatchedBy(func(ip shared_domain.Ip) bool {
			return ip.String() == "102.165.35.0"
		})).Return(pkg.Some(&geo_domain.Country{ISOCode: "FR"}), nil)
		geoRepo.On("FindByIp", mock.Anything, mock.MatchedBy(func(ip shared_domain.Ip) bool {
			return ip.String() == "100.42.176.0"
		})).Return(pkg.Some(&geo_domain.Country{ISOCode: "DE"}), nil)
		repo.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.LinkCountryViewCounter](), nil)
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkCountryViewCounter")).Return(nil)
		handler := NewIncremeterOnLinkVisitProcessed(logger, geoRepo, repo)
		linkId := "019af069-4ef2-7068-9dbb-48f7e1b9dd37"
		event := domain_link_visit.NewLinkVisitsProcessedDomainEvent(
			linkId,
			8,
			4,
			[]domain_link_visit.IpsVisits{
				{Ip: "102.177.191.0", TotalViews: 2, UniqueViews: 1},
				{Ip: "102.212.56.0", TotalViews: 1, UniqueViews: 1},
				{Ip: "102.129.65.0", TotalViews: 2, UniqueViews: 2},
				{Ip: "102.165.35.0", TotalViews: 1, UniqueViews: 1},
				{Ip: "100.42.176.0", TotalViews: 2, UniqueViews: 2},
			},
		)
		err := handler.On(context.Background(), event)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		calls := 0
		for _, country := range []string{"ES", "FR", "DE"} {
			countryCode, _ := shared_domain.NewCountryCode(country)
			repo.AssertCalled(t, "Save", mock.Anything, mock.MatchedBy(func(counter domain.LinkCountryViewCounter) bool {
				return counter.CountryCode == countryCode
			}))
			calls++
		}
		if calls != 3 {
			t.Errorf("Expected 3 Save calls, got %d", calls)
		}
	})
}