package application

import (
	"context"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
)

type ClientCreator struct {
	logger domain_shared.Logger
	repo   domain.ClientRepository
}

func NewClientCreator(logger domain_shared.Logger, repo domain.ClientRepository) *ClientCreator {
	return &ClientCreator{
		logger: logger,
		repo:   repo,
	}
}

func (c ClientCreator) Run(ctx context.Context, name, email string) (*domain.Client, error) {
	c.logger.Info(ctx, "ClientCreator - Run - Params into: ", domain_shared.NewField("name", name))
	client, err := domain.NewClient(name, email)
	if err != nil {
		c.logger.Error(ctx, "ClientCreator - Run - Error creating client", domain_shared.NewField("error", err.Error()))
		return nil, err
	}

	err = c.repo.Save(ctx, *client)
	if err != nil {
		c.logger.Error(ctx, "ClientCreator - Run - Error saving client", domain_shared.NewField("error", err.Error()))
		return nil, err
	}

	return client, nil
}
