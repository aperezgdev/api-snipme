package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (pkg.Optional[*RefreshToken], error)
	FindByUserId(ctx context.Context, userId shared_domain.Id) ([]*RefreshToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserId(ctx context.Context, userId shared_domain.Id) error
}

type RefreshTokenRepositoryMock struct {
	mock.Mock
}

func (m *RefreshTokenRepositoryMock) Save(ctx context.Context, token *RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *RefreshTokenRepositoryMock) FindByToken(ctx context.Context, token string) (pkg.Optional[*RefreshToken], error) {
	args := m.Called(ctx, token)
	return args.Get(0).(pkg.Optional[*RefreshToken]), args.Error(1)
}

func (m *RefreshTokenRepositoryMock) FindByUserId(ctx context.Context, userId shared_domain.Id) ([]*RefreshToken, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]*RefreshToken), args.Error(1)
}

func (m *RefreshTokenRepositoryMock) Delete(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *RefreshTokenRepositoryMock) DeleteByUserId(ctx context.Context, userId shared_domain.Id) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
