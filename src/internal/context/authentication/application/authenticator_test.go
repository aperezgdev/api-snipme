package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticatorAuthenticateWithOAuth(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("creates new user when not exists", func(t *testing.T) {
		t.Parallel()
		userRepo := &domain.UserRepositoryMock{}
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)

		provider := domain.OAuthProviderGoogle
		oauthSubject := "google-user-123"
		email := "test@example.com"

		userRepo.On("FindByOAuthProviderAndSubject", mock.Anything, provider, oauthSubject).
			Return(pkg.EmptyOptional[*domain.User](), nil)
		userRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		jwtManager.On("Generate", mock.Anything, email).Return("jwt-token", nil)
		refreshTokenRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		result, err := authenticator.Run(context.Background(), provider, oauthSubject, email)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if result.AccessToken != "jwt-token" {
			t.Errorf("Expected access token 'jwt-token', got %s", result.AccessToken)
		}
		if result.RefreshToken == "" {
			t.Error("Expected refresh token, got empty")
		}
		if result.ExpiresIn != 1200 {
			t.Errorf("Expected ExpiresIn 1200, got %d", result.ExpiresIn)
		}
		if string(result.User.Email) != email {
			t.Errorf("Expected email %s, got %s", email, result.User.Email)
		}

		userRepo.AssertExpectations(t)
		refreshTokenRepo.AssertExpectations(t)
		jwtManager.AssertExpectations(t)
		eventBus.AssertExpectations(t)
	})

	t.Run("updates existing user email if changed", func(t *testing.T) {
		t.Parallel()
		userRepo := &domain.UserRepositoryMock{}
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)

		provider := domain.OAuthProviderGoogle
		oauthSubject := "google-user-123"
		newEmail := "newemail@example.com"

		existingUser, _ := domain.NewUser("oldemail@example.com", provider, oauthSubject)

		userRepo.On("FindByOAuthProviderAndSubject", mock.Anything, provider, oauthSubject).
			Return(pkg.Some(existingUser), nil)
		userRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
		userRepo.On("FindById", mock.Anything, mock.Anything).Return(pkg.Some(existingUser), nil)

		jwtManager.On("Generate", mock.Anything, newEmail).Return("jwt-token", nil)
		refreshTokenRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

		result, err := authenticator.Run(context.Background(), provider, oauthSubject, newEmail)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(result.User.Email) != newEmail {
			t.Errorf("Expected email %s, got %s", newEmail, result.User.Email)
		}

		userRepo.AssertExpectations(t)
		refreshTokenRepo.AssertExpectations(t)
		jwtManager.AssertExpectations(t)
	})

	t.Run("fails when user save returns error", func(t *testing.T) {
		t.Parallel()
		userRepo := &domain.UserRepositoryMock{}
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)

		provider := domain.OAuthProviderGoogle
		oauthSubject := "google-user-123"
		email := "test@example.com"

		userRepo.On("FindByOAuthProviderAndSubject", mock.Anything, provider, oauthSubject).
			Return(pkg.EmptyOptional[*domain.User](), nil)
		userRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(errors.New("database error"))

		_, err := authenticator.Run(context.Background(), provider, oauthSubject, email)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		userRepo.AssertExpectations(t)
	})

	t.Run("fails when JWT generation returns error", func(t *testing.T) {
		t.Parallel()
		userRepo := &domain.UserRepositoryMock{}
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}
		eventBus := &shared_domain.EventBusMock{}

		authenticator := NewAuthenticator(logger, userRepo, refreshTokenRepo, jwtManager, eventBus, 30, 20)

		provider := domain.OAuthProviderGoogle
		oauthSubject := "google-user-123"
		email := "test@example.com"

		userRepo.On("FindByOAuthProviderAndSubject", mock.Anything, provider, oauthSubject).
			Return(pkg.EmptyOptional[*domain.User](), nil)
		userRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
		eventBus.On("Publish", mock.Anything, mock.Anything).Return()

		jwtManager.On("Generate", mock.Anything, email).Return("", errors.New("jwt error"))

		_, err := authenticator.Run(context.Background(), provider, oauthSubject, email)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		userRepo.AssertExpectations(t)
		jwtManager.AssertExpectations(t)
	})
}
