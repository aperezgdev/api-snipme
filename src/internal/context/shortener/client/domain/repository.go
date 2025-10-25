package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

type ClientRepository interface {
	Save(ctx context.Context, client Client) error
	Remove(ctx context.Context, id shared_domain.Id) error
}

type ShortLinkRepository struct {
	mock.Mock
}

func (m *ShortLinkRepository) Save(ctx context.Context, client Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *ShortLinkRepository) Remove(ctx context.Context, id shared_domain.Id) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
