package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type ClientRepository interface {
	FindById(ctx context.Context, id shared_domain.Id) (pkg.Optional[*Client], error)
	Save(ctx context.Context, client Client) error
	Remove(ctx context.Context, id shared_domain.Id) error
}

type ClientRepositoryMock struct {
	mock.Mock
}

func (m *ClientRepositoryMock) FindById(ctx context.Context, id shared_domain.Id) (pkg.Optional[*Client], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(pkg.Optional[*Client]), args.Error(1)
}

func (m *ClientRepositoryMock) Save(ctx context.Context, client Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *ClientRepositoryMock) Remove(ctx context.Context, id shared_domain.Id) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
