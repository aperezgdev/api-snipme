package infrastructure

import (
	"context"
	"time"

	"net/netip"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcLinkVisitRepository struct {
	logger  shared_domain_context.Logger
	queries *generated.Queries
}

func NewSqlcLinkVisitRepository(logger shared_domain_context.Logger, q *generated.Queries) *SqlcLinkVisitRepository {
	return &SqlcLinkVisitRepository{queries: q, logger: logger}
}

func (r *SqlcLinkVisitRepository) Save(ctx context.Context, linkVisit domain.LinkVisit) error {
	r.logger.Info(ctx, "SqlcLinkVisitRepository - Save - Params into",
		shared_domain_context.NewField("linkVisitId", linkVisit.Id.String()),
		shared_domain_context.NewField("linkId", linkVisit.LinkId.String()),
		shared_domain_context.NewField("ip", linkVisit.Ip),
	)

	id := pgtype.UUID{}
	_ = id.Scan(linkVisit.Id.String())

	createdOn := pgtype.Timestamptz{}
	createdOn.Time = time.Time(linkVisit.CreatedOn)
	createdOn.Valid = true

	linkVisitId := pgtype.UUID{}
	_ = linkVisitId.Scan(linkVisit.LinkId.String())

	addr := netip.AddrPort(linkVisit.Ip).Addr()

	params := generated.SaveLinkVisitParams{
		ID:        id,
		LinkID:    linkVisitId,
		Ip:        &addr,
		UserAgent: pgtype.Text{String: string(linkVisit.UserAgent), Valid: true},
		CreatedOn: createdOn,
	}
	err := r.queries.SaveLinkVisit(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkVisitRepository - Save - Error saving link visit", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkVisitRepository - Save - Link visit saved successfully", shared_domain_context.NewField("linkVisitId", linkVisit.Id.String()))
	return nil
}
