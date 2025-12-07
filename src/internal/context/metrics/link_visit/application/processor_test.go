package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkVisitProcessor_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success with visits", func(t *testing.T) {
		repo := &domain.LinkVisitRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		processor := NewLinkVisitProcessor(logger, repo, eventBus)

		linkId, _ := shared_domain_context.ParseID("11111111-1111-1111-1111-111111111111")
		visitor1, _ := shared_domain_context.ParseID("22222222-2222-2222-2222-222222222222")
		visitor2, _ := shared_domain_context.ParseID("33333333-3333-3333-3333-333333333333")
		ip1, _ := shared_domain.NewIp("1.1.1.1")
		ip2, _ := shared_domain.NewIp("2.2.2.2")

		visits := []domain.LinkVisit{
			{LinkId: linkId, VisitorId: visitor1, Ip: ip1},
			{LinkId: linkId, VisitorId: visitor2, Ip: ip1},
			{LinkId: linkId, VisitorId: visitor1, Ip: ip2},
		}

		repo.On("FindOlds", mock.Anything).Return(visits, nil)
		repo.On("RemoveAll", mock.Anything, mock.Anything).Return(nil)

		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		err := processor.Run(context.Background())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		repo.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})

	t.Run("Run with no visits", func(t *testing.T) {
		repo := &domain.LinkVisitRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		processor := NewLinkVisitProcessor(logger, repo, eventBus)

		repo.On("FindOlds", mock.Anything).Return([]domain.LinkVisit{}, nil)

		err := processor.Run(context.Background())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run fails on finder error", func(t *testing.T) {
		repo := &domain.LinkVisitRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		processor := NewLinkVisitProcessor(logger, repo, eventBus)

		repo.On("FindOlds", mock.Anything).Return([]domain.LinkVisit{}, errors.New("db error"))

		err := processor.Run(context.Background())
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})

	t.Run("Run fails on remover error", func(t *testing.T) {
		repo := &domain.LinkVisitRepositoryMock{}
		eventBus := &shared_domain_context.EventBusMock{}
		processor := NewLinkVisitProcessor(logger, repo, eventBus)

		linkId, _ := shared_domain_context.ParseID("11111111-1111-1111-1111-111111111111")
		visitor1, _ := shared_domain_context.ParseID("22222222-2222-2222-2222-222222222222")
		ip1, _ := shared_domain.NewIp("1.1.1.1")
		visits := []domain.LinkVisit{{LinkId: linkId, VisitorId: visitor1, Ip: ip1}}

		repo.On("FindOlds", mock.Anything).Return(visits, nil)
		repo.On("RemoveAll", mock.Anything, mock.Anything).Return(errors.New("remove error"))

		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		err := processor.Run(context.Background())
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		repo.AssertExpectations(t)
	})
}

func TestLinkVisitProcessor_calculateVisits(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}
	repo := &domain.LinkVisitRepositoryMock{}
	eventBus := &shared_domain_context.EventBusMock{}
	processor := NewLinkVisitProcessor(logger, repo, eventBus)

	linkId, _ := shared_domain_context.ParseID("11111111-1111-1111-1111-111111111111")
	visitor1, _ := shared_domain_context.ParseID("22222222-2222-2222-2222-222222222222")
	visitor2, _ := shared_domain_context.ParseID("33333333-3333-3333-3333-333333333333")
	ip1, _ := shared_domain.NewIp("1.1.1.1")
	ip2, _ := shared_domain.NewIp("2.2.2.2")

	visits := []domain.LinkVisit{
		{LinkId: linkId, VisitorId: visitor1, Ip: ip1},
		{LinkId: linkId, VisitorId: visitor1, Ip: ip1},
		{LinkId: linkId, VisitorId: visitor2, Ip: ip1},
		{LinkId: linkId, VisitorId: visitor1, Ip: ip2},
	}
	grouped := map[shared_domain_context.Id][]domain.LinkVisit{
		linkId: visits,
	}
	events := processor.calculateVisits(context.Background(), grouped)
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	processed, ok := events[0].(domain.LinkVisitsProcessed)
	if !ok {
		t.Fatalf("Expected event type LinkVisitsProcessed")
	}
	if processed.TotalViews != 4 {
		t.Errorf("Expected TotalViews 3, got %d", processed.TotalViews)
	}
	if processed.UniqueViews != 2 {
		t.Errorf("Expected UniqueViews 2, got %d", processed.UniqueViews)
	}

	want := map[string]struct {
		total  uint
		unique uint
	}{
		"1.1.1.1": {3, 2},
		"2.2.2.2": {1, 1},
	}
	for _, ipf := range processed.IpsFrequency {
		w, ok := want[ipf.Ip]
		if !ok {
			t.Errorf("Unexpected IP %s", ipf.Ip)
		}
		if ipf.TotalViews != w.total {
			t.Errorf("For IP %s, expected TotalViews %d, got %d", ipf.Ip, w.total, ipf.TotalViews)
		}
		if ipf.UniqueViews != w.unique {
			t.Errorf("For IP %s, expected UniqueViews %d, got %d", ipf.Ip, w.unique, ipf.UniqueViews)
		}
	}
}
