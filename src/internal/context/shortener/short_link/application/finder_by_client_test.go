package application

import (
	"context"
	"errors"
	"testing"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	domain_client "github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/stretchr/testify/mock"
)

func TestFinderByClient(t *testing.T) {
	logger := shared_domain_context.DummyLogger{}
	t.Run("Run successfully finds short links by client ID", func(t *testing.T) {
		repoMock := &domain.ShortLinkRepositoryMock{}
		clientRepoMock := &domain_client.ClientRepositoryMock{}
		finder := NewShortLinkFinderByClient(logger, repoMock, clientRepoMock)

		clientId := "00000000-0000-0000-0000-000000000001"
		idVO, _ := shared_domain_context.ParseID(clientId)
		expectedLinks := []*domain.ShortLink{
			{
				OriginalRoute: domain.ShortLinkOriginalRoute("http://example.com/1"),
				Client:        idVO,
			},
			{
				OriginalRoute: domain.ShortLinkOriginalRoute("http://example.com/2"),
				Client:        idVO,
			},
		}

		repoMock.On("FindByClient", mock.Anything, idVO).Return(expectedLinks, nil)
		clientRepoMock.On("FindById", mock.Anything, idVO).Return(
			pkg.Some(&domain_client.Client{
				Id: idVO,
			}), nil)

		links, err := finder.Run(context.Background(), clientId)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(links) != len(expectedLinks) {
			t.Fatalf("expected %d links, got %d", len(expectedLinks), len(links))
		}
		repoMock.AssertExpectations(t)
		clientRepoMock.AssertExpectations(t)
	})

	t.Run("Returns error when client not found", func(t *testing.T) {
		repoMock := &domain.ShortLinkRepositoryMock{}
		clientRepoMock := &domain_client.ClientRepositoryMock{}
		finder := NewShortLinkFinderByClient(logger, repoMock, clientRepoMock)

		clientId := "00000000-0000-0000-0000-000000000002"
		idVO, _ := shared_domain_context.ParseID(clientId)

		clientRepoMock.On("FindById", mock.Anything, idVO).Return(
			pkg.EmptyOptional[*domain_client.Client](), nil)

		_, err := finder.Run(context.Background(), clientId)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
		var notFoundErr shared_domain_context.NotFoundError
		if !errors.As(err, &notFoundErr) {
			t.Fatalf("expected NotFoundError, got %v", err)
		}

		clientRepoMock.AssertExpectations(t)
	})
}
