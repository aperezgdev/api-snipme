package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/infrastructure"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestTokenRefresherRefresh(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("refreshes token successfully", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		userRepo := &domain.UserRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}

		refresher := NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)

		user, _ := domain.NewUser("test@example.com", domain.OAuthProviderGoogle, "subject-123")
		refreshToken, _ := domain.NewRefreshToken(user.Id.String(), "refresh-token-string", time.Now().Add(24*time.Hour))

		refreshTokenRepo.On("FindByToken", mock.Anything, "refresh-token-string").
			Return(pkg.Some(refreshToken), nil)
		userRepo.On("FindById", mock.Anything, user.Id).Return(pkg.Some(user), nil)
		jwtManager.On("Generate", user.Id.String(), "test@example.com").Return("new-jwt-token", nil)

		result, err := refresher.Refresh(context.Background(), "refresh-token-string")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected result, got nil")
		}
		if result.AccessToken != "new-jwt-token" {
			t.Errorf("Expected access token 'new-jwt-token', got %s", result.AccessToken)
		}
		if result.ExpiresIn != 1200 {
			t.Errorf("Expected ExpiresIn 1200, got %d", result.ExpiresIn)
		}

		refreshTokenRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		jwtManager.AssertExpectations(t)
	})

	t.Run("fails when refresh token not found", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		userRepo := &domain.UserRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}

		refresher := NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)

		refreshTokenRepo.On("FindByToken", mock.Anything, "invalid-token").
			Return(pkg.EmptyOptional[*domain.RefreshToken](), nil)

		_, err := refresher.Refresh(context.Background(), "invalid-token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		refreshTokenRepo.AssertExpectations(t)
	})

	t.Run("fails when refresh token is expired", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		userRepo := &domain.UserRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}

		refresher := NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)

		userId, _ := shared_domain.NewID()
		tokenId, _ := shared_domain.NewID()
		expiredToken, _ := domain.RefreshTokenFromPrimitives(
			tokenId.String(),
			userId.String(),
			"expired-token",
			time.Now().Add(-24*time.Hour),
			time.Now().Add(-48*time.Hour),
		)

		refreshTokenRepo.On("FindByToken", mock.Anything, "expired-token").
			Return(pkg.Some(expiredToken), nil)

		_, err := refresher.Refresh(context.Background(), "expired-token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		refreshTokenRepo.AssertExpectations(t)
	})

	t.Run("fails when user not found", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		userRepo := &domain.UserRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}

		refresher := NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)

		userId, _ := shared_domain.NewID()
		refreshToken, _ := domain.NewRefreshToken(userId.String(), "refresh-token-string", time.Now().Add(24*time.Hour))

		refreshTokenRepo.On("FindByToken", mock.Anything, "refresh-token-string").
			Return(pkg.Some(refreshToken), nil)
		userRepo.On("FindById", mock.Anything, userId).Return(pkg.EmptyOptional[*domain.User](), nil)

		_, err := refresher.Refresh(context.Background(), "refresh-token-string")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		refreshTokenRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("fails when JWT generation returns error", func(t *testing.T) {
		t.Parallel()
		refreshTokenRepo := &domain.RefreshTokenRepositoryMock{}
		userRepo := &domain.UserRepositoryMock{}
		jwtManager := &infrastructure.JWTManagerMock{}

		refresher := NewTokenRefresher(logger, refreshTokenRepo, userRepo, jwtManager, 20)

		user, _ := domain.NewUser("test@example.com", domain.OAuthProviderGoogle, "subject-123")
		refreshToken, _ := domain.NewRefreshToken(user.Id.String(), "refresh-token-string", time.Now().Add(24*time.Hour))

		refreshTokenRepo.On("FindByToken", mock.Anything, "refresh-token-string").
			Return(pkg.Some(refreshToken), nil)
		userRepo.On("FindById", mock.Anything, mock.Anything).Return(pkg.Some(user), nil)
		jwtManager.On("Generate", mock.Anything, "test@example.com").
			Return("", errors.New("jwt generation failed"))

		_, err := refresher.Refresh(context.Background(), "refresh-token-string")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		refreshTokenRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
		jwtManager.AssertExpectations(t)
	})
}
