package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	domain_client "github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetShortLinkByClientHTTPHandler_Handler(t *testing.T) {
	t.Parallel()

	logger := shared_domain_context.DummyLogger{}

	t.Run("Successfully returns short links for a client", func(t *testing.T) {
		t.Parallel()

		client, _ := domain_client.NewClient("Test Client", "test@example.com")
		shortLink1, _ := domain.NewShortLink("Summary 1", "https://example.com/1", client.Id.String())
		shortLink2, _ := domain.NewShortLink("Summary 2", "https://example.com/2", client.Id.String())

		mockClientRepo := new(domain_client.ClientRepositoryMock)
		mockClientRepo.On("FindById", context.Background(), client.Id).Return(pkg.Some(client), nil)

		mockShortLinkRepo := new(domain.ShortLinkRepositoryMock)
		mockShortLinkRepo.On("FindByClient", context.Background(), client.Id).Return([]*domain.ShortLink{shortLink1, shortLink2}, nil)

		finder := application.NewShortLinkFinderByClient(logger, mockShortLinkRepo, mockClientRepo)
		handler := NewGetShortLinkByClientHTTPHandler(logger, *finder)

		req := httptest.NewRequest(http.MethodGet, "/short-link?client_id="+client.Id.String(), nil)
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []shortLinkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 short links, got %d", len(response))
		}

		mockClientRepo.AssertExpectations(t)
		mockShortLinkRepo.AssertExpectations(t)
	})

	t.Run("Returns not found when client does not exist", func(t *testing.T) {
		t.Parallel()

		clientId, _ := shared_domain_context.NewID()

		mockClientRepo := new(domain_client.ClientRepositoryMock)
		mockClientRepo.On("FindById", context.Background(), clientId).Return(pkg.EmptyOptional[*domain_client.Client](), nil)

		mockShortLinkRepo := new(domain.ShortLinkRepositoryMock)

		finder := application.NewShortLinkFinderByClient(logger, mockShortLinkRepo, mockClientRepo)
		handler := NewGetShortLinkByClientHTTPHandler(logger, *finder)

		req := httptest.NewRequest(http.MethodGet, "/short-link?client_id="+clientId.String(), nil)
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

		mockClientRepo.AssertExpectations(t)
	})

	t.Run("Returns internal server error on repository error", func(t *testing.T) {
		t.Parallel()

		clientId, _ := shared_domain_context.NewID()

		mockClientRepo := new(domain_client.ClientRepositoryMock)
		mockClientRepo.On("FindById", context.Background(), clientId).Return(pkg.EmptyOptional[*domain_client.Client](), assert.AnError)

		mockShortLinkRepo := new(domain.ShortLinkRepositoryMock)

		finder := application.NewShortLinkFinderByClient(logger, mockShortLinkRepo, mockClientRepo)
		handler := NewGetShortLinkByClientHTTPHandler(logger, *finder)

		req := httptest.NewRequest(http.MethodGet, "/short-link?client_id="+clientId.String(), nil)
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		mockClientRepo.AssertExpectations(t)
	})
}
