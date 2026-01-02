package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	refresh_token_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/domain"
        refresh_token_application "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/application"
	refresh_token_infrastructure "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/infrastructure"
	user_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func TestPostRefreshTokenHandler(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("fails with empty refresh token", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}

		refresher := refresh_token_application.NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)
		handler := NewPostRefreshTokenHandler(logger, refresher)

		body := map[string]string{"refresh_token": ""}
		b, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("fails with invalid JSON body", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}

		refresher := refresh_token_application.NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)
		handler := NewPostRefreshTokenHandler(logger, refresher)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Handler(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("has correct route", func(t *testing.T) {
		refreshTokenRepo := &refresh_token_domain.RefreshTokenRepositoryMock{}
		userRepo := &user_domain.UserRepositoryMock{}
		jwtManager := &refresh_token_infrastructure.JWTManagerMock{}

		refresher := refresh_token_application.NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)
		handler := NewPostRefreshTokenHandler(logger, refresher)

		route := handler.Route()
		if route != "/auth/refresh" {
			t.Errorf("Expected route /auth/refresh, got %s", route)
		}
	})
}
