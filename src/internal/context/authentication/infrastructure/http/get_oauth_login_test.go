package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/stretchr/testify/mock"
)

func TestGetOAuthLoginHandler(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("redirects to OAuth provider successfully", func(t *testing.T) {
		t.Parallel()
		oauthClient := &infrastructure.OAuthClientMock{}

		handler := NewGetOAuthLoginHandler(logger, oauthClient, "google", "test-secret-key")

		oauthClient.On("GetAuthURL", mock.Anything).Return("https://oauth.example.com/authorize?state=random")

		req := httptest.NewRequest(http.MethodGet, "/auth/google/login", nil)
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusTemporaryRedirect && resp.StatusCode != http.StatusFound {
			t.Errorf("Expected status 302 or 307, got %d", resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location == "" {
			t.Error("Expected Location header, got empty")
		}

		oauthClient.AssertExpectations(t)
	})

	t.Run("has correct route", func(t *testing.T) {
		oauthClient := &infrastructure.OAuthClientMock{}
		handler := NewGetOAuthLoginHandler(logger, oauthClient, "google", "test-secret-key")

		route := handler.Route()
		if route != "/auth/google/login" {
			t.Errorf("Expected route /auth/google/login, got %s", route)
		}
	})
}
