package application

import (
	"context"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	domain_client "github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
)

type ShortLinkFinderByClient struct {
	logger       shared_domain_context.Logger
	repo         domain.ShortLinkRepository
	clientFinder domain_client.ClientFinder
}

func NewShortLinkFinderByClient(
	logger shared_domain_context.Logger,
	repo domain.ShortLinkRepository,
	userRepo domain_client.ClientRepository,
) *ShortLinkFinderByClient {
	return &ShortLinkFinderByClient{
		logger:       logger,
		repo:         repo,
		clientFinder: *domain_client.NewClientFinder(logger, userRepo),
	}
}

func (f *ShortLinkFinderByClient) Run(ctx context.Context, clientId string) ([]*domain.ShortLink, error) {
	f.logger.Info(ctx, "FinderByClient - Run = Params into:", shared_domain_context.NewField("clientId", clientId))

	_, err := f.clientFinder.Run(ctx, clientId)
	if err != nil {
		f.logger.Error(ctx, "FinderByClient - Run = Error finding client:", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	idVO, err := shared_domain_context.ParseID(clientId)
	if err != nil {
		f.logger.Error(ctx, "FinderByClient - Run = Error creating client ID VO:", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	links, err := f.repo.FindByClient(ctx, idVO)
	if err != nil {
		f.logger.Error(ctx, "FinderByClient - Run = Error finding short links by client:", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	f.logger.Info(ctx, "FinderByClient - Run = Successfully found short links by client:", shared_domain_context.NewField("count", len(links)))
	return links, nil
}
