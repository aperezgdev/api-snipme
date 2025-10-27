package domain

import (
	"context"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type ShortLinkRepository interface {
	Save(ctx context.Context, shortLink *ShortLink) error
	FindByCode(ctx context.Context, code ShortLinkCode) (pkg.Optional[*ShortLink], error)
	Remove(ctx context.Context, id domain_shared.Id) error
}

type ShortLinkRepositoryMock struct {
	mock.Mock
}

func (m *ShortLinkRepositoryMock) Save(ctx context.Context, shortLink *ShortLink) error {
	args := m.Called(ctx, shortLink)
	return args.Error(0)
}

func (m *ShortLinkRepositoryMock) FindByCode(ctx context.Context, code ShortLinkCode) (pkg.Optional[*ShortLink], error) {
	args := m.Called(ctx, code)
	return args.Get(0).(pkg.Optional[*ShortLink]), args.Error(1)
}

func (m *ShortLinkRepositoryMock) Remove(ctx context.Context, id domain_shared.Id) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
