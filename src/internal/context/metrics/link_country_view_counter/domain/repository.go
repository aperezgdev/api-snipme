package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type LinkCountryViewCounterRepository interface {
	Find(ctx context.Context, idLink shared_domain.Id, countryCode CountryCode) (pkg.Optional[LinkCountryViewCounter], error)
	Save(ctx context.Context, linkCountryStats LinkCountryViewCounter) error
	Update(ctx context.Context, linkCountryStats LinkCountryViewCounter) error
	RemoveByLink(ctx context.Context, idLink shared_domain.Id) error
}

type LinkCountryViewCounterRepositoryMock struct {
	mock.Mock
}

func (m *LinkCountryViewCounterRepositoryMock) Find(ctx context.Context, idLink shared_domain.Id, countryCode CountryCode) (pkg.Optional[LinkCountryViewCounter], error) {
	args := m.Called(ctx, idLink, countryCode)
	return args.Get(0).(pkg.Optional[LinkCountryViewCounter]), args.Error(1)
}

func (m *LinkCountryViewCounterRepositoryMock) Save(ctx context.Context, linkCountryStats LinkCountryViewCounter) error {
	args := m.Called(ctx, linkCountryStats)
	return args.Error(0)
}

func (m *LinkCountryViewCounterRepositoryMock) Update(ctx context.Context, linkCountryStats LinkCountryViewCounter) error {
	args := m.Called(ctx, linkCountryStats)
	return args.Error(0)
}

func (m *LinkCountryViewCounterRepositoryMock) RemoveByLink(ctx context.Context, idLink shared_domain.Id) error {
	args := m.Called(ctx, idLink)
	return args.Error(0)
}
