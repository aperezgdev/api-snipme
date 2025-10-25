package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	link_analytics_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcLinkAnalyticsRepository struct {
	queries *generated.Queries
}

func NewSqlcLinkAnalyticsRepository(q *generated.Queries) *SqlcLinkAnalyticsRepository {
	return &SqlcLinkAnalyticsRepository{queries: q}
}

func (r *SqlcLinkAnalyticsRepository) Save(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	id := pgtype.UUID{}
	_ = id.Scan(analytics.Id.String())

	linkId := pgtype.UUID{}
	_ = linkId.Scan(analytics.LinkId.String())

	createdOn := pgtype.Timestamptz{}
	createdOn.Time = time.Time(analytics.UpdateOn)
	createdOn.Valid = true

	params := generated.SaveLinkAnalyticsParams{
		ID:             id,
		LinkID:         linkId,
		TotalViews:     pgtype.Int4{Int32: int32(analytics.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(analytics.UniqueViews), Valid: true},
		CreatedOn:      createdOn,
	}
	return r.queries.SaveLinkAnalytics(ctx, params)
}

func (r *SqlcLinkAnalyticsRepository) Update(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	id := pgtype.UUID{}
	_ = id.Scan(analytics.Id.String())

	params := generated.UpdateLinkAnalyticsParams{
		ID:             id,
		TotalViews:     pgtype.Int4{Int32: int32(analytics.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(analytics.UniqueViews), Valid: true},
		CreatedOn:      pgtype.Timestamptz{Time: time.Time(analytics.UpdateOn), Valid: true},
	}
	return r.queries.UpdateLinkAnalytics(ctx, params)
}

func (r *SqlcLinkAnalyticsRepository) RemoveByLink(ctx context.Context, idLink shared_domain_context.Id) error {
	linkId := pgtype.UUID{}
	_ = linkId.Scan(idLink.String())
	return r.queries.RemoveLinkAnalyticsByLink(ctx, linkId)
}
