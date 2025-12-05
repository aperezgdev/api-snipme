package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	link_analytics_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcLinkAnalyticsRepository struct {
	logger  shared_domain_context.Logger
	queries *generated.Queries
}

func NewSqlcLinkAnalyticsRepository(logger shared_domain_context.Logger, q *generated.Queries) *SqlcLinkAnalyticsRepository {
	return &SqlcLinkAnalyticsRepository{queries: q, logger: logger}
}

func (r *SqlcLinkAnalyticsRepository) Save(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Save - Params into", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
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
	err := r.queries.SaveLinkAnalytics(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - Save - Error at save", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Save - Save successfully", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
	return nil
}

func (r *SqlcLinkAnalyticsRepository) FindByLinkId(ctx context.Context, linkId shared_domain_context.Id) (pkg.Optional[*link_analytics_domain.LinkAnalytics], error) {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - FindByLinkId - Params into", shared_domain_context.NewField("linkId", linkId.String()))
	id := pgtype.UUID{}
	_ = id.Scan(linkId.String())

	rows, err := r.queries.FindByLinkID(ctx, id)
	if rows == nil || len(rows) == 0 {
		r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - FindByLinkId - LinkAnalytics not found", shared_domain_context.NewField("linkId", linkId.String()))
		return pkg.EmptyOptional[*link_analytics_domain.LinkAnalytics](), nil
	}

	if err != nil {
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - FindByLinkId - Error at finding", shared_domain_context.NewField("error", err.Error()))
		return pkg.EmptyOptional[*link_analytics_domain.LinkAnalytics](), err
	}

	row := rows[0]

	linkAnalytics := link_analytics_domain.LinkAnalytics{
		Id:          shared_domain_context.Id(row.ID.Bytes),
		LinkId:      shared_domain_context.Id(row.LinkID.Bytes),
		TotalViews:  shared_domain.LinkViewsCounter(row.TotalViews.Int32),
		UniqueViews: shared_domain.LinkViewsCounter(row.UniqueVisitors.Int32),
		UpdateOn:    shared_domain_context.UpdatedOn(row.CreatedOn.Time),
	}

	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - FindByLinkId - Find successfully", shared_domain_context.NewField("linkId", linkId.String()))
	return pkg.Some(&linkAnalytics), nil
}

func (r *SqlcLinkAnalyticsRepository) Update(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Update - Params into", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
	id := pgtype.UUID{}
	_ = id.Scan(analytics.Id.String())

	params := generated.UpdateLinkAnalyticsParams{
		ID:             id,
		TotalViews:     pgtype.Int4{Int32: int32(analytics.TotalViews), Valid: true},
		UniqueVisitors: pgtype.Int4{Int32: int32(analytics.UniqueViews), Valid: true},
		CreatedOn:      pgtype.Timestamptz{Time: time.Time(analytics.UpdateOn), Valid: true},
	}
	err := r.queries.UpdateLinkAnalytics(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - Update - Error trying to update", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Update - Update successfully", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
	return nil
}

func (r *SqlcLinkAnalyticsRepository) RemoveByLink(ctx context.Context, idLink shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Params into", shared_domain_context.NewField("idLink", idLink.String()))
	linkId := pgtype.UUID{}
	_ = linkId.Scan(idLink.String())
	err := r.queries.RemoveLinkAnalyticsByLink(ctx, linkId)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Error at delete", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Delete successfully", shared_domain_context.NewField("idLink", idLink.String()))
	return nil
}
