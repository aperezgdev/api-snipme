package application

import (
	"context"
	"errors"
	"net/netip"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestLinkVisitCreator_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid link visit data", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkVisitRepositoryMock{}
		creator := NewLinkVisitCreator(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		visitorId := "11111111-1111-1111-1111-111111111111"
		ip := "192.168.1.1"
		userAgent := "Mozilla/5.0"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkVisit")).Return(nil)

		result, err := creator.Run(context.Background(), linkId, visitorId, ip, userAgent)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if netip.Addr(result.Ip).String() != ip || string(result.UserAgent) != userAgent {
			t.Errorf("Expected link visit with ip %s and userAgent %s, got %s and %s", ip, userAgent, netip.Addr(result.Ip).String(), string(result.UserAgent))
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid link visit data", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkVisitRepositoryMock{}
		creator := NewLinkVisitCreator(logger, repo)

		linkId := ""
		visitorId := ""
		ip := ""
		userAgent := ""

		_, err := creator.Run(context.Background(), linkId, visitorId, ip, userAgent)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Run fails on repository save error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.LinkVisitRepositoryMock{}
		creator := NewLinkVisitCreator(logger, repo)

		linkId := "00000000-0000-0000-0000-000000000000"
		visitorId := "11111111-1111-1111-1111-111111111111"
		ip := "192.168.1.1"
		userAgent := "Mozilla/5.0"

		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.LinkVisit")).Return(errors.New("database error"))

		_, err := creator.Run(context.Background(), linkId, visitorId, ip, userAgent)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
