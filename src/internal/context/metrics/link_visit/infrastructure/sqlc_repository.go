package infrastructure

import (
	"context"
	"time"

	"net/netip"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcLinkVisitRepository struct {
	queries *generated.Queries
}

func NewSqlcLinkVisitRepository(q *generated.Queries) *SqlcLinkVisitRepository {
	return &SqlcLinkVisitRepository{queries: q}
}

func (r *SqlcLinkVisitRepository) Save(ctx context.Context, linkVisit domain.LinkVisit) error {
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
	return r.queries.SaveLinkVisit(ctx, params)
}
