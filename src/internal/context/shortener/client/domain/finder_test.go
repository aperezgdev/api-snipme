package domain

import (
	"context"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestClientFinder(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}
	t.Run("FindById returns client when found", func(t *testing.T) {
		repo := &ClientRepositoryMock{}
		finder := NewClientFinder(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		client := Client{
			Id: domainId,
		}

		repo.On("FindById", mock.Anything, domainId).Return(pkg.Some(&client), nil)

		foundClient, err := finder.Run(context.Background(), id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if foundClient.Id != client.Id {
			t.Fatalf("Expected client ID %v, got %v", client.Id, foundClient.Id)
		}

		repo.AssertExpectations(t)
	})

	t.Run("FindById returns empty when not found", func(t *testing.T) {
		repo := &ClientRepositoryMock{}
		finder := NewClientFinder(logger, repo)

		id := "00000000-0000-0000-0000-000000000000"
		domainId, _ := shared_domain_context.ParseID(id)

		repo.On("FindById", mock.Anything, domainId).Return(pkg.EmptyOptional[*Client](), nil)

		_, err := finder.Run(context.Background(), id)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		repo.AssertExpectations(t)
	})
}
