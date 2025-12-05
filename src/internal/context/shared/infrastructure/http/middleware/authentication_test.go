package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/application"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

type mockRoute struct {
	handlerFunc http.HandlerFunc
	route       string
	method      string
}

func (m *mockRoute) Handler(w http.ResponseWriter, r *http.Request) {
	m.handlerFunc(w, r)
}

func (m *mockRoute) Route() string {
	return m.route
}

func (m *mockRoute) Method() string {
	return m.method
}

func TestAuthenticationMiddleware(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("allows request with valid token", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := application.NewTokenValidator(logger, jwtManager, userRepo)
		middleware := NewAuthenticationMiddleware(logger, validator)

		user, _ := domain.NewUser(
			"test@example.com",
			domain.OAuthProviderGoogle,
			"oauth-subject-123",
		)

		claims := &domain.TokenClaims{
			UserID: user.Id.String(),
			Email:  "test@example.com",
		}

		jwtManager.On("Validate", "valid-token").Return(claims, nil)
		userRepo.On("FindById", mock.Anything, user.Id).Return(pkg.Some(user), nil)

		handlerCalled := false
		nextHandler := &mockRoute{
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true

				authenticatedUser, ok := GetAuthenticatedUser(r.Context())
				if !ok {
					t.Error("Expected authenticated user in context")
				}
				if authenticatedUser.Email != "test@example.com" {
					t.Errorf("Expected email test@example.com, got %s", authenticatedUser.Email)
				}

				w.WriteHeader(http.StatusOK)
			},
			route:  "/test",
			method: http.MethodGet,
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := middleware.Handle(nextHandler)
		handler(w, req)

		if !handlerCalled {
			t.Error("Next handler was not called")
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		jwtManager.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("rejects request without Authorization header", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := application.NewTokenValidator(logger, jwtManager, userRepo)
		middleware := NewAuthenticationMiddleware(logger, validator)

		handlerCalled := false
		nextHandler := &mockRoute{
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			},
			route:  "/test",
			method: http.MethodGet,
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		handler := middleware.Handle(nextHandler)
		handler(w, req)

		if handlerCalled {
			t.Error("Next handler should not be called")
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects request with invalid Authorization header format", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := application.NewTokenValidator(logger, jwtManager, userRepo)
		middleware := NewAuthenticationMiddleware(logger, validator)

		handlerCalled := false
		nextHandler := &mockRoute{
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			},
			route:  "/test",
			method: http.MethodGet,
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		w := httptest.NewRecorder()

		handler := middleware.Handle(nextHandler)
		handler(w, req)

		if handlerCalled {
			t.Error("Next handler should not be called")
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("rejects request with invalid token", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := application.NewTokenValidator(logger, jwtManager, userRepo)
		middleware := NewAuthenticationMiddleware(logger, validator)

		jwtManager.On("Validate", "invalid-token").Return(nil, errors.New("invalid token"))

		handlerCalled := false
		nextHandler := &mockRoute{
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
			},
			route:  "/test",
			method: http.MethodGet,
		}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler := middleware.Handle(nextHandler)
		handler(w, req)

		if handlerCalled {
			t.Error("Next handler should not be called")
		}

		resp := w.Result()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}

		jwtManager.AssertExpectations(t)
	})
}

func TestGetAuthenticatedUser(t *testing.T) {
	t.Run("returns user when present in context", func(t *testing.T) {
		user, _ := domain.NewUser(
			"test@example.com",
			domain.OAuthProviderGoogle,
			"oauth-subject-123",
		)

		ctx := context.WithValue(context.Background(), AuthenticatedUserKey, user)

		retrievedUser, ok := GetAuthenticatedUser(ctx)
		if !ok {
			t.Error("Expected user to be present in context")
		}

		if retrievedUser.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", retrievedUser.Email)
		}
	})

	t.Run("returns false when user not in context", func(t *testing.T) {
		ctx := context.Background()

		_, ok := GetAuthenticatedUser(ctx)
		if ok {
			t.Error("Expected user to not be present in context")
		}
	})
}
