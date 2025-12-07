package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type CountryRepository interface {
	FindByIp(ctx context.Context, ip shared_domain.Ip) (pkg.Optional[*Country], error)
}

type CountryRepositoryMock struct {
	mock.Mock
}

func (m *CountryRepositoryMock) FindByIp(ctx context.Context, ip shared_domain.Ip) (pkg.Optional[*Country], error) {
	args := m.Called(ctx, ip)
	return args.Get(0).(pkg.Optional[*Country]), args.Error(1)
}
