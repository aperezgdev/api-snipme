package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteShortLinkHTTPHandler_Handler(t *testing.T) {
	t.Parallel()

	logger := shared_domain_context.DummyLogger{}

	t.Run("Successfully deletes a short link", func(t *testing.T) {
		shortLink, _ := domain.NewPublicShortLink("https://example.com")
		mockRepo := new(domain.ShortLinkRepositoryMock)
		mockRepo.On("FindById", mock.Anything, mock.MatchedBy(func(id shared_domain_context.Id) bool {
			return id == shortLink.Id
		})).Return(pkg.Some(shortLink), nil)
		mockRepo.On("Remove", mock.Anything, mock.MatchedBy(func(id shared_domain_context.Id) bool {
			return id == shortLink.Id
		})).Return(nil)

		remover := application.NewShortLinkRemover(logger, mockRepo)
		handler := NewDeleteShortLinkHTTPHandler(logger, *remover)

		req := httptest.NewRequest(http.MethodDelete, "/short-link/"+shortLink.Id.String(), nil)
		req.SetPathValue("short_code", shortLink.Id.String())
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Returns not found when short link does not exist", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		mockRepo.On("FindById", mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.ShortLink](), nil)

		remover := application.NewShortLinkRemover(logger, mockRepo)
		handler := NewDeleteShortLinkHTTPHandler(logger, *remover)

		nonExistentId, _ := shared_domain_context.NewID()
		req := httptest.NewRequest(http.MethodDelete, "/short-link/"+nonExistentId.String(), nil)
		req.SetPathValue("short_code", nonExistentId.String())
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Returns internal server error on repository error", func(t *testing.T) {
		mockRepo := new(domain.ShortLinkRepositoryMock)
		mockRepo.On("FindById", mock.Anything, mock.Anything).Return(pkg.EmptyOptional[*domain.ShortLink](), assert.AnError)

		remover := application.NewShortLinkRemover(logger, mockRepo)
		handler := NewDeleteShortLinkHTTPHandler(logger, *remover)

		id, _ := shared_domain_context.NewID()
		req := httptest.NewRequest(http.MethodDelete, "/short-link/"+id.String(), nil)
		req.SetPathValue("short_code", id.String())
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		mockRepo.AssertExpectations(t)
	})
}
