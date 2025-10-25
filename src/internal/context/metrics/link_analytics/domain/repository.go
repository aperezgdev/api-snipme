package domain

import (
	"context"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

type LinkAnalyticsRepository interface {
	Save(ctx context.Context, linkAnalytics LinkAnalytics) error
	Update(ctx context.Context, linkAnalytics LinkAnalytics) error
	RemoveByLink(ctx context.Context, idLink shared_domain.Id) error
}

type LinkAnalyticsRepositoryMock struct {
	mock.Mock
}

func (m *LinkAnalyticsRepositoryMock) Save(ctx context.Context, linkAnalytics LinkAnalytics) error {
	args := m.Called(ctx, linkAnalytics)
	return args.Error(0)
}

func (m *LinkAnalyticsRepositoryMock) Update(ctx context.Context, linkAnalytics LinkAnalytics) error {
	args := m.Called(ctx, linkAnalytics)
	return args.Error(0)
}

func (m *LinkAnalyticsRepositoryMock) RemoveByLink(ctx context.Context, idLink shared_domain.Id) error {
	args := m.Called(ctx, idLink)
	return args.Error(0)
}
