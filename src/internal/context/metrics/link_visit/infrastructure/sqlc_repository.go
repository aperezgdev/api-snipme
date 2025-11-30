package infrastructure

import (
	"context"
	"time"

	"net/netip"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
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

	addr := netip.Addr(linkVisit.Ip)

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

func (r *SqlcLinkVisitRepository) FindOlds(ctx context.Context) ([]domain.LinkVisit, error) {
	r.logger.Info(ctx, "SqlcLinkVisitRepository - FindOld - Finding old link visits")
	rows, err := r.queries.FindOldLinkVisits(ctx)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkVisitRepository - FindOld - Error finding old link visits", shared_domain_context.NewField("error", err.Error()))
		return nil, err
	}

	linkVisits := pkg.Map(rows, func(row generated.LinkVisit) domain.LinkVisit {
		id, _ := shared_domain_context.ParseID(row.ID.String())
		linkId, _ := shared_domain_context.ParseID(row.LinkID.String())

		return domain.LinkVisit{
			Id:        id,
			LinkId:    linkId,
			Ip:        domain.LinkVisitIP(*row.Ip),
			UserAgent: domain.LinkVisitUserAgent(row.UserAgent.String),
			CreatedOn: shared_domain_context.CreatedOn(row.CreatedOn.Time),
		}
	})
	r.logger.Info(ctx, "SqlcLinkVisitRepository - FindOld - Found old link visits", shared_domain_context.NewField("count", len(linkVisits)))
	return linkVisits, nil
}

func (r *SqlcLinkVisitRepository) RemoveAll(ctx context.Context, linkVisitsIds []shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcLinkVisitRepository - RemoveAll - Removing link visits", shared_domain_context.NewField("count", len(linkVisitsIds)))

	var ids []pgtype.UUID
	for _, id := range linkVisitsIds {
		pgId := pgtype.UUID{}
		_ = pgId.Scan(id.String())
		ids = append(ids, pgId)
	}

	err := r.queries.RemoveLinkVisits(ctx, ids)
	if err != nil {
		r.logger.Error(ctx, "SqlcLinkVisitRepository - RemoveAll - Error removing link visits", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcLinkVisitRepository - RemoveAll - Link visits removed successfully", shared_domain_context.NewField("count", len(linkVisitsIds)))
	return nil
}
