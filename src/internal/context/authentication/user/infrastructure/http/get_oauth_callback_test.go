package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	refresh_token_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/domain"
	refresh_token_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/infrastructure"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/application"
	user_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestGetOAuthCallbackHandler(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("handles callback successfully", func(t *testing.T) {
		t.Parallel()
		oauthClient := &infrastructure.OAuthClientMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := application.NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)
		stateSecret := "test-secret-key"
		handler := NewGetOAuthCallbackHandler(logger, oauthClient, authenticator, user_domain.OAuthProviderGoogle, stateSecret)

		state := "test-state-value"
		h := hmac.New(sha256.New, []byte(stateSecret))
		h.Write([]byte(state))
		signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
		stateCookie := state + "." + signature

		token := &oauth2.Token{AccessToken: "access-token"}
		userInfo := &infrastructure.OAuthUserInfo{
			Email:   "test@example.com",
			Subject: "oauth-subject-123",
		}

		oauthClient.On("ExchangeCode", mock.Anything, "auth-code").Return(token, nil)
		oauthClient.On("GetUserInfo", mock.Anything, token).Return(userInfo, nil)

		userRepo.On("FindByOAuthProviderAndSubject", mock.Anything, user_domain.OAuthProviderGoogle, "oauth-subject-123").
			Return(pkg.EmptyOptional[*user_domain.User](), nil)
		userRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
		jwtManager.On("Generate", mock.Anything, "test@example.com").Return("jwt-token", nil)
		refreshTokenRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=auth-code&state="+state, nil)
		req.AddCookie(&http.Cookie{Name: "oauth_state", Value: stateCookie})
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		oauthClient.AssertExpectations(t)
	})

	t.Run("fails when code parameter is missing", func(t *testing.T) {
		t.Parallel()
		oauthClient := &infrastructure.OAuthClientMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := application.NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)
		handler := NewGetOAuthCallbackHandler(logger, oauthClient, authenticator, user_domain.OAuthProviderGoogle, "test-secret-key")

		req := httptest.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("has correct route", func(t *testing.T) {
		oauthClient := &infrastructure.OAuthClientMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := application.NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)
		handler := NewGetOAuthCallbackHandler(logger, oauthClient, authenticator, user_domain.OAuthProviderGoogle, "test-secret-key")

		route := handler.Route()
		if route != "/auth/google/callback" {
			t.Errorf("Expected route /auth/google/callback, got %s", route)
		}
	})
}
