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
	logger  shared_domain_context.Logger
}

func NewSqlcLinkAnalyticsRepository(q *generated.Queries, logger shared_domain_context.Logger) *SqlcLinkAnalyticsRepository {
	return &SqlcLinkAnalyticsRepository{queries: q, logger: logger}
}

func (r *SqlcLinkAnalyticsRepository) Save(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Save - Entrada", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
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
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - Save - Error al guardar", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Save - Guardado exitosamente", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
	return nil
}

func (r *SqlcLinkAnalyticsRepository) Update(ctx context.Context, analytics link_analytics_domain.LinkAnalytics) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Update - Entrada", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
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
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - Update - Error al actualizar", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - Update - Actualizado exitosamente", shared_domain_context.NewField("linkAnalyticsId", analytics.Id.String()))
	return nil
}

func (r *SqlcLinkAnalyticsRepository) RemoveByLink(ctx context.Context, idLink shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Entrada", shared_domain_context.NewField("idLink", idLink.String()))
	linkId := pgtype.UUID{}
	_ = linkId.Scan(idLink.String())
	err := r.queries.RemoveLinkAnalyticsByLink(ctx, linkId)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Error al eliminar", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkAnalyticsRepository - RemoveByLink - Eliminado exitosamente", shared_domain_context.NewField("idLink", idLink.String()))
	return nil
}
