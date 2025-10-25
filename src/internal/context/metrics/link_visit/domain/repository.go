package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LinkVisitRepository interface {
	Save(ctx context.Context, linkVisit LinkVisit) error
}

type LinkVisitRepositoryMock struct {
	mock.Mock
}

func (m *LinkVisitRepositoryMock) Save(ctx context.Context, linkVisit LinkVisit) error {
	args := m.Called(ctx, linkVisit)
	return args.Error(0)
}
