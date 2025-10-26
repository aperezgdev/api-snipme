package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	domain_shared "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcLinkCountryViewCounterRepository struct {
	queries *generated.Queries
	logger  shared_domain_context.Logger
}

func NewSqlcLinkCountryViewCounterRepository(q *generated.Queries, logger shared_domain_context.Logger) *SqlcLinkCountryViewCounterRepository {
	return &SqlcLinkCountryViewCounterRepository{queries: q, logger: logger}
}

func (r *SqlcLinkCountryViewCounterRepository) Find(ctx context.Context, idLink shared_domain_context.Id, countryCode domain.CountryCode) (pkg.Optional[domain.LinkCountryViewCounter], error) {
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Find - Params into", shared_domain_context.NewField("idLink", idLink.String()), shared_domain_context.NewField("countryCode", string(countryCode)))
	linkID := pgtype.UUID{}
	_ = linkID.Scan(idLink.String())

	params := generated.FindLinkCounterViewCounterParams{
		LinkID:      linkID,
		CountryCode: string(countryCode),
	}

	result, err := r.queries.FindLinkCounterViewCounter(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkCountryViewCounterRepository - Find - Error querying", shared_domain_context.NewField("error", err.Error()))
		return pkg.Optional[domain.LinkCountryViewCounter]{}, err
	}

	entity := domain.LinkCountryViewCounter{
		Id:          shared_domain_context.Id(result.ID.Bytes),
		LinkId:      shared_domain_context.Id(result.LinkID.Bytes),
		CountryCode: domain.CountryCode(result.CountryCode),
		TotalViews:  domain_shared.LinkViewsCounter(result.TotalViews.Int32),
		UniqueViews: domain_shared.LinkViewsCounter(result.UniqueVisitors.Int32),
		CreatedOn:   shared_domain_context.CreatedOn(result.CreatedOn.Time),
	}

	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Find - Success", shared_domain_context.NewField("id", entity.Id.String()))
	return pkg.Some(entity), nil
}

func (r *SqlcLinkCountryViewCounterRepository) Save(ctx context.Context, linkCountryStats domain.LinkCountryViewCounter) error {
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Save - Params into", shared_domain_context.NewField("id", linkCountryStats.Id.String()), shared_domain_context.NewField("linkId", linkCountryStats.LinkId.String()), shared_domain_context.NewField("countryCode", string(linkCountryStats.CountryCode)))
	id := pgtype.UUID{}
	_ = id.Scan(linkCountryStats.Id.String())

	linkID := pgtype.UUID{}
	_ = linkID.Scan(linkCountryStats.LinkId.String())

	createdOn := pgtype.Timestamptz{}
	createdOn.Time = time.Time(linkCountryStats.CreatedOn)
	createdOn.Valid = true

	params := generated.SaveLinkCounterViewCounterParams{
		ID:             id,
		LinkID:         linkID,
		CountryCode:    string(linkCountryStats.CountryCode),
		TotalViews:     pgtype.Int4{Int32: int32(linkCountryStats.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(linkCountryStats.UniqueViews), Valid: true},
		CreatedOn:      createdOn,
	}
	err := r.queries.SaveLinkCounterViewCounter(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkCountryViewCounterRepository - Save - Error saving", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Save - Success", shared_domain_context.NewField("id", linkCountryStats.Id.String()))
	return nil
}

func (r *SqlcLinkCountryViewCounterRepository) Update(ctx context.Context, linkCountryStats domain.LinkCountryViewCounter) error {
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Update - Params into", shared_domain_context.NewField("linkId", linkCountryStats.LinkId.String()), shared_domain_context.NewField("countryCode", string(linkCountryStats.CountryCode)))
	linkID := pgtype.UUID{}
	_ = linkID.Scan(linkCountryStats.LinkId.String())

	params := generated.UpdateLinkCounterViewCounterParams{
		TotalViews:     pgtype.Int4{Int32: int32(linkCountryStats.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(linkCountryStats.UniqueViews), Valid: true},
		LinkID:         linkID,
		CountryCode:    string(linkCountryStats.CountryCode),
	}
	err := r.queries.UpdateLinkCounterViewCounter(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkCountryViewCounterRepository - Update - Error updating", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - Update - Success", shared_domain_context.NewField("linkId", linkID.String()))
	return nil
}

func (r *SqlcLinkCountryViewCounterRepository) RemoveByLink(ctx context.Context, idLink shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - RemoveByLink - Params into", shared_domain_context.NewField("idLink", idLink.String()))
	linkID := pgtype.UUID{}
	_ = linkID.Scan(idLink.String())
	err := r.queries.RemoveLinkCounterViewCounterByLink(ctx, linkID)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkCountryViewCounterRepository - RemoveByLink - Error removing", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkCountryViewCounterRepository - RemoveByLink - Success", shared_domain_context.NewField("idLink", idLink.String()))
	return nil
}
