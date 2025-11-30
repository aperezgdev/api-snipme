package domain

import (
	"context"

	domain_shared_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

type LinkVisitRepository interface {
	Save(ctx context.Context, linkVisit LinkVisit) error
	FindOlds(ctx context.Context) ([]LinkVisit, error)
	RemoveAll(ctx context.Context, linkVisitsIds []domain_shared_context.Id) error
}

type LinkVisitRepositoryMock struct {
	mock.Mock
}

func (m *LinkVisitRepositoryMock) Save(ctx context.Context, linkVisit LinkVisit) error {
	args := m.Called(ctx, linkVisit)
	return args.Error(0)
}

func (m *LinkVisitRepositoryMock) FindOlds(ctx context.Context) ([]LinkVisit, error) {
	args := m.Called(ctx)
	return args.Get(0).([]LinkVisit), args.Error(1)
}

func (m *LinkVisitRepositoryMock) RemoveAll(ctx context.Context, linkVisitsIds []domain_shared_context.Id) error {
	args := m.Called(ctx, linkVisitsIds)
	return args.Error(0)
}
