package domain

import (
	"context"
	"testing"

	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestUserEmailUpdater(t *testing.T) {
	logger := shared_domain.DummyLogger{}

	t.Run("successfully updates user email", func(t *testing.T) {
		t.Parallel()
		userRepo := &UserRepositoryMock{}

		emailUpdater := NewUserEmailUpdater(logger, userRepo)

		existingUser, _ := NewUser(
			"joe@doe.com",
			OAuthProviderGoogle,
			"google-user-123",
		)

		newEmail := "joe@doe.com"
		userRepo.On("FindById", mock.Anything, mock.Anything).
			Return(pkg.Some(existingUser), nil)
		userRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		err := emailUpdater.Run(context.Background(), existingUser.Id.String(), newEmail)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		userRepo.AssertExpectations(t)
	})

	t.Run("fails when user not found", func(t *testing.T) {
		t.Parallel()
		userRepo := &UserRepositoryMock{}

		emailUpdater := NewUserEmailUpdater(logger, userRepo)

		userId, _ := shared_domain.NewID()
		newEmail := "joe@doe.com"
		userRepo.On("FindById", mock.Anything, userId).
			Return(pkg.EmptyOptional[*User](), nil)

		err := emailUpdater.Run(context.Background(), userId.String(), newEmail)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		userRepo.AssertExpectations(t)
	})
}