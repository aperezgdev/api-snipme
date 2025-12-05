package application

import (
	"context"
	"errors"
	"testing"

	auth_domain "github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/stretchr/testify/mock"
)

func TestCreatorOnUserCreated_On(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("On creates client when UserCreated event is received", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewCreatorOnUserCreated(logger, repo)

		userId := "550e8400-e29b-41d4-a716-446655440000"
		email := "newuser@example.com"
		name := "New User"

		event := auth_domain.NewUserCreatedEvent(userId, email, name)
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.Client")).Return(nil)

		err := creator.On(context.Background(), event)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		repo.AssertExpectations(t)
		repo.AssertCalled(t, "Save", mock.Anything, mock.AnythingOfType("domain.Client"))
	})

	t.Run("On handles repository error gracefully", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewCreatorOnUserCreated(logger, repo)

		userId := "550e8400-e29b-41d4-a716-446655440000"
		email := "newuser@example.com"
		name := "New User"

		event := auth_domain.NewUserCreatedEvent(userId, email, name)
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.Client")).Return(errors.New("database error"))

		err := creator.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})

	t.Run("On handles invalid client data from event", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewCreatorOnUserCreated(logger, repo)

		userId := "550e8400-e29b-41d4-a716-446655440000"
		email := "invalid-email"
		name := ""

		event := auth_domain.NewUserCreatedEvent(userId, email, name)

		err := creator.On(context.Background(), event)
		if err == nil {
			t.Fatalf("Expected error for invalid client data, got nil")
		}

		repo.AssertNotCalled(t, "Save")
	})

	t.Run("SubscribedTo returns UserCreated event", func(t *testing.T) {
		repo := &domain.ClientRepositoryMock{}
		creator := NewCreatorOnUserCreated(logger, repo)

		events := creator.SubscribedTo()
		if len(events) != 1 {
			t.Fatalf("Expected 1 subscribed event, got %d", len(events))
		}
		if events[0] != auth_domain.UserCreatedEventName {
			t.Errorf("Expected %s, got %s", auth_domain.UserCreatedEventName, events[0])
		}
	})
}
