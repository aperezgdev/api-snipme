package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

type LinkViewsCounterRemover struct {
	logger     shared_domain.Logger
	repository domain.LinkCountryViewCounterRepository
}

func NewLinkViewsCounterRemover(logger shared_domain.Logger, repository domain.LinkCountryViewCounterRepository) *LinkViewsCounterRemover {
	return &LinkViewsCounterRemover{
		logger:     logger,
		repository: repository,
	}
}

func (r LinkViewsCounterRemover) Run(ctx context.Context, linkId string) error {
	r.logger.Info(ctx, "LinkViewsCounterRemover - Run - Params into: ", shared_domain.NewField("linkId", linkId))
	domainId, err := shared_domain.ParseID(linkId)
	if err != nil {
		r.logger.Error(ctx, "LinkViewsCounterRemover - Run - Error parsing linkId", shared_domain.NewField("error", err.Error()))
		return err
	}
	err = r.repository.RemoveByLink(ctx, domainId)
	if err != nil {
		r.logger.Error(ctx, "LinkViewsCounterRemover - Run - Error removing link views counter", shared_domain.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "LinkViewsCounterRemover - Run - Successfully removed link views counter", shared_domain.NewField("linkId", linkId))
	return nil
}
