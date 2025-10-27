package application

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/stretchr/testify/mock"
)

func TestClientCreator_Run(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}

	t.Run("Run success on valid client data", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewClientCreator(logger, repo)

		name := "John Doe"
		email := "john@example.com"
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.Client")).Return(nil)

		result, err := creator.Run(context.Background(), name, email)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(result.Name) != name || string(result.Email) != email {
			t.Errorf("Expected client with name %s and email %s, got %s and %s", name, email, result.Name, result.Email)
		}

		repo.AssertExpectations(t)
	})

	t.Run("Run fails on invalid client data", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewClientCreator(logger, repo)

		name := ""
		email := "invalid-email"

		_, err := creator.Run(context.Background(), name, email)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})

	t.Run("Run fails on repository save error", func(t *testing.T) {
		t.Parallel()
		repo := &domain.ClientRepositoryMock{}
		creator := NewClientCreator(logger, repo)

		name := "Jane Doe"
		email := "jane@example.com"
		repo.On("Save", mock.Anything, mock.AnythingOfType("domain.Client")).Return(errors.New("database error"))

		_, err := creator.Run(context.Background(), name, email)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
