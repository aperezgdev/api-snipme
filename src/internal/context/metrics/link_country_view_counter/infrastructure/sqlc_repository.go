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
}

func NewSqlcLinkCountryViewCounterRepository(q *generated.Queries) *SqlcLinkCountryViewCounterRepository {
	return &SqlcLinkCountryViewCounterRepository{queries: q}
}

func (r *SqlcLinkCountryViewCounterRepository) Find(ctx context.Context, idLink shared_domain_context.Id, countryCode domain.CountryCode) (pkg.Optional[domain.LinkCountryViewCounter], error) {
	linkID := pgtype.UUID{}
	_ = linkID.Scan(idLink.String())

	params := generated.FindLinkCounterViewCounterParams{
		LinkID:      linkID,
		CountryCode: string(countryCode),
	}

	result, err := r.queries.FindLinkCounterViewCounter(ctx, params)
	if err != nil {
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

	return pkg.Some(entity), nil
}

func (r *SqlcLinkCountryViewCounterRepository) Save(ctx context.Context, linkCountryStats domain.LinkCountryViewCounter) error {
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
	return r.queries.SaveLinkCounterViewCounter(ctx, params)
}

func (r *SqlcLinkCountryViewCounterRepository) Update(ctx context.Context, linkCountryStats domain.LinkCountryViewCounter) error {
	linkID := pgtype.UUID{}
	_ = linkID.Scan(linkCountryStats.LinkId.String())

	params := generated.UpdateLinkCounterViewCounterParams{
		TotalViews:     pgtype.Int4{Int32: int32(linkCountryStats.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(linkCountryStats.UniqueViews), Valid: true},
		LinkID:         linkID,
		CountryCode:    string(linkCountryStats.CountryCode),
	}
	return r.queries.UpdateLinkCounterViewCounter(ctx, params)
}

func (r *SqlcLinkCountryViewCounterRepository) RemoveByLink(ctx context.Context, idLink shared_domain_context.Id) error {
	linkID := pgtype.UUID{}
	_ = linkID.Scan(idLink.String())
	return r.queries.RemoveLinkCounterViewCounterByLink(ctx, linkID)
}
