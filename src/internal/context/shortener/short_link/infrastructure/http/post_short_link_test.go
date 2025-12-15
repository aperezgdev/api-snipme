package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostShortLinkHTTPHandler_Handler(t *testing.T) {
	t.Parallel()

	logger := shared_domain_context.DummyLogger{}
	eventBus := shared_domain_context.NewEventBusInMemory()
	eventBusPtr := &eventBus

	t.Run("Successfully creates a short link", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		mockRepo.On("Save", context.Background(), mock.AnythingOfType("*domain.ShortLink")).Return(nil)

		creator := application.NewShortLinkCreator(logger, mockRepo, eventBusPtr)
		handler := NewPostShortLinkHTTPHandler(logger, *creator)

		body := map[string]string{
			"original_url": "https://example.com",
			"summary":      "Test",
			"client_id":    "",
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/short-links", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Returns bad request on invalid JSON", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		creator := application.NewShortLinkCreator(logger, mockRepo, eventBusPtr)
		handler := NewPostShortLinkHTTPHandler(logger, *creator)

		req := httptest.NewRequest(http.MethodPost, "/short-links", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Returns bad request on validation error", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		creator := application.NewShortLinkCreator(logger, mockRepo, eventBusPtr)
		handler := NewPostShortLinkHTTPHandler(logger, *creator)

		body := map[string]string{
			"original_url": "invalid-url",
			"summary":      "Test",
			"client_id":    "",
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/short-links", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Returns internal server error on repository error", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		mockRepo.On("Save", context.Background(), mock.AnythingOfType("*domain.ShortLink")).Return(assert.AnError)

		creator := application.NewShortLinkCreator(logger, mockRepo, eventBusPtr)
		handler := NewPostShortLinkHTTPHandler(logger, *creator)

		body := map[string]string{
			"original_url": "https://example.com",
			"summary":      "Test",
			"client_id":    "",
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/short-links", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		mockRepo.AssertExpectations(t)
	})
}
