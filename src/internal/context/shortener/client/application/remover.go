package application

import (
	"context"

	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
)

type ClientRemover struct {
	logger domain_shared.Logger
	repo   domain.ClientRepository
}

func NewClientRemover(logger domain_shared.Logger, repo domain.ClientRepository) *ClientRemover {
	return &ClientRemover{
		logger: logger,
		repo:   repo,
	}
}

func (r ClientRemover) Run(ctx context.Context, id string) error {
	r.logger.Info(ctx, "ClientRemover - Run - Params into: ", domain_shared.NewField("id", id))
	domainId, err := domain_shared.ParseID(id)
	if err != nil {
		r.logger.Error(ctx, "ClientRemover - Run - Error parsing id", domain_shared.NewField("error", err.Error()))
		return err
	}
	err = r.repo.Remove(ctx, domainId)
	if err != nil {
		r.logger.Error(ctx, "ClientRemover - Run - Error removing client", domain_shared.NewField("error", err.Error()))
		return err
	}
	return nil
}
