package domain

import (
	"context"

	shared_context_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
)

type LinkCountryViewCounterFinder struct {
	logger shared_context_domain.Logger
	repo   LinkCountryViewCounterRepository
}

func NewLinkCountryViewCounterFinder(
	logger shared_context_domain.Logger,
	repo LinkCountryViewCounterRepository,
) LinkCountryViewCounterFinder {
	return LinkCountryViewCounterFinder{
		logger: logger,
		repo:   repo,
	}
}

func (f LinkCountryViewCounterFinder) Run(ctx context.Context, idLink, countryCode string) (pkg.Optional[LinkCountryViewCounter], error) {
	f.logger.Info(ctx, "LinkCountryViewCounterFinder - Run - Params into: ", shared_context_domain.NewField("idLink", idLink))
	domainId, err := shared_context_domain.ParseID(idLink)
	if err != nil {
		f.logger.Error(ctx, "LinkCountryViewCounterFinder - Run - Error parsing idLink", shared_context_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[LinkCountryViewCounter](), err
	}

	countryCodeVO, err := NewCountryCode(countryCode)
	if err != nil {
		f.logger.Error(ctx, "LinkCountryViewCounterFinder - Run - Error creating country code VO", shared_context_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[LinkCountryViewCounter](), err
	}

	linkCountryViewCounter, err := f.repo.Find(ctx, domainId, countryCodeVO)
	if err != nil {
		f.logger.Error(ctx, "LinkCountryViewCounterFinder - Run - Error finding link country view counter", shared_context_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[LinkCountryViewCounter](), err
	}
	f.logger.Info(ctx, "LinkCountryViewCounterFinder - Run - Successfully found link country view counter", shared_context_domain.NewField("idLink", idLink))

	return linkCountryViewCounter, nil
}
