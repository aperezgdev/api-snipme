package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindById(ctx context.Context, id shared_domain.Id) (pkg.Optional[*User], error)
	FindByEmail(ctx context.Context, email string) (pkg.Optional[*User], error)
	FindByOAuthProviderAndSubject(ctx context.Context, provider OAuthProvider, subject string) (pkg.Optional[*User], error)
	Update(ctx context.Context, user *User) error
}

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Save(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindById(ctx context.Context, id shared_domain.Id) (pkg.Optional[*User], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(pkg.Optional[*User]), args.Error(1)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (pkg.Optional[*User], error) {
	args := m.Called(ctx, email)
	return args.Get(0).(pkg.Optional[*User]), args.Error(1)
}

func (m *UserRepositoryMock) FindByOAuthProviderAndSubject(ctx context.Context, provider OAuthProvider, subject string) (pkg.Optional[*User], error) {
	args := m.Called(ctx, provider, subject)
	return args.Get(0).(pkg.Optional[*User]), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
