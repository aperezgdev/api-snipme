package application

import (
	"context"
	"errors"
	"testing"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/user/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/infrastructure"
        refresh_token_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestTokenValidator_Validate(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("validates token successfully", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := NewTokenValidator(logger, jwtManager, userRepo)

		user, _ := domain.NewUser("test@example.com", domain.OAuthProviderGoogle, "subject-123")

		claims := &refresh_token_domain.TokenClaims{
			UserID: user.Id.String(),
			Email:  "test@example.com",
		}

		jwtManager.On("Validate", "valid-token").Return(claims, nil)
		userRepo.On("FindById", mock.Anything, user.Id).Return(pkg.Some(user), nil)

		result, err := validator.Validate(context.Background(), "valid-token")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected user, got nil")
		}
		if result.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", result.Email)
		}

		jwtManager.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("fails on invalid token", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := NewTokenValidator(logger, jwtManager, userRepo)

		jwtManager.On("Validate", "invalid-token").Return(nil, errors.New("invalid token"))

		_, err := validator.Validate(context.Background(), "invalid-token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		jwtManager.AssertExpectations(t)
	})

	t.Run("fails when user not found", func(t *testing.T) {
		t.Parallel()
		jwtManager := &infrastructure.JWTManagerMock{}
		userRepo := &domain.UserRepositoryMock{}

		validator := NewTokenValidator(logger, jwtManager, userRepo)

		userId, _ := shared_domain.NewID()
		claims := &refresh_token_domain.TokenClaims{
			UserID: userId.String(),
			Email:  "test@example.com",
		}

		jwtManager.On("Validate", "valid-token").Return(claims, nil)
		userRepo.On("FindById", mock.Anything, userId).Return(pkg.EmptyOptional[*domain.User](), nil)

		_, err := validator.Validate(context.Background(), "valid-token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		jwtManager.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}
