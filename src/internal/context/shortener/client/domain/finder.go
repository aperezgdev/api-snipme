package domain

import (
	"context"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type ClientFinder struct {
	logger shared_domain_context.Logger
	repo   ClientRepository
}

func NewClientFinder(
	logger shared_domain_context.Logger,
	repo ClientRepository,
) *ClientFinder {
	return &ClientFinder{
		logger: logger,
		repo:   repo,
	}
}

func (f *ClientFinder) Run(ctx context.Context, clientId string) (*Client, error) {
	f.logger.Info(ctx, "ClientFinder - Run = Params into:", shared_domain_context.NewField("clientId", clientId))

	idVO, err := shared_domain_context.ParseID(clientId)
	if err != nil {
		f.logger.Error(ctx, "ClientFinder - Run = Error creating client ID VO:", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	clientOpt, err := f.repo.FindById(ctx, idVO)
	if err != nil {
		f.logger.Error(ctx, "ClientFinder - Run = Error finding client by ID:", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	if !clientOpt.IsPresent() {
		f.logger.Error(ctx, "ClientFinder - Run = Client not found:", shared_domain_context.NewField("clientId", clientId))
		return nil, domain_shared.NewNotFoundError("Client not found")
	}

	client := clientOpt.Get()
	f.logger.Info(ctx, "ClientFinder - Run = Successfully found client by ID:", shared_domain_context.NewField("clientId", client.Id.String()))
	return client, nil
}
