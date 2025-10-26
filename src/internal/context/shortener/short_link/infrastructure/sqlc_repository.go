package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcShortLinkRepository struct {
	logger  shared_domain_context.Logger
	queries *generated.Queries
}

func NewSqlcShortLinkRepository(logger shared_domain_context.Logger, q *generated.Queries) *SqlcShortLinkRepository {
	return &SqlcShortLinkRepository{queries: q, logger: logger}
}

func (r *SqlcShortLinkRepository) Save(ctx context.Context, shortLink *domain.ShortLink) error {
	r.logger.Info(ctx, "SqlcShortLinkRepository - Save - Params into", shared_domain_context.NewField("id", shortLink.Id.String()), shared_domain_context.NewField("code", string(shortLink.Code)))
	id := pgtype.UUID{}
	_ = id.Scan(shortLink.Id.String())

	clientID := pgtype.UUID{}
	if shortLink.Client.String() != "" && shortLink.Client.String() != "00000000-0000-0000-0000-000000000000" {
		_ = clientID.Scan(shortLink.Client.String())
		clientID.Valid = true
	} else {
		clientID.Valid = false
	}

	createdOn := pgtype.Timestamptz{}
	createdOn.Time = time.Time(shortLink.CreatedOn)
	createdOn.Valid = true

	err := r.queries.SaveShortLink(ctx, generated.SaveShortLinkParams{
		ID:            id,
		OriginalRoute: string(shortLink.OriginalRoute),
		ClientID:      clientID,
		Code:          string(shortLink.Code),
		CreatedOn:     createdOn,
	})

	if err != nil {
		r.logger.Error(ctx, "SqlcShortLinkRepository - Save - Error saving short link", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcShortLinkRepository - Save - Short link saved successfully", shared_domain_context.NewField("id", shortLink.Id.String()))
	return nil
}

func (r *SqlcShortLinkRepository) Remove(ctx context.Context, id shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcShortLinkRepository - Remove - Params into", shared_domain_context.NewField("id", id.String()))
	uuid := pgtype.UUID{}
	_ = uuid.Scan(id.String())
	err := r.queries.RemoveShortLink(ctx, uuid)
	if err != nil {
		r.logger.Error(ctx, "SqlcShortLinkRepository - Remove - Error removing short link", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcShortLinkRepository - Remove - Short link removed successfully", shared_domain_context.NewField("id", id.String()))
	return nil
}

func (r *SqlcShortLinkRepository) FindByCode(ctx context.Context, code domain.ShortLinkCode) (pkg.Optional[*domain.ShortLink], error) {
	r.logger.Info(ctx, "SqlcShortLinkRepository - FindByCode - Params into", shared_domain_context.NewField("code", string(code)))
	shortLink, err := r.queries.FindShortLinkByCode(ctx, string(code))
	if err != nil {
		r.logger.Error(ctx, "SqlcShortLinkRepository - FindByCode - Error finding short link", shared_domain_context.NewField("error", err.Error()))
		return pkg.Optional[*domain.ShortLink]{}, err
	}

	id, _ := shared_domain_context.ParseID(shortLink.ID.String())
	clientID, _ := shared_domain_context.ParseID(shortLink.ClientID.String())
	result := &domain.ShortLink{
		Id:            id,
		OriginalRoute: domain.ShortLinkOriginalRoute(shortLink.OriginalRoute),
		Code:          domain.ShortLinkCode(shortLink.Code),
		Client:        clientID,
		CreatedOn:     shared_domain_context.CreatedOn(shortLink.CreatedOn.Time),
	}
	r.logger.Info(ctx, "SqlcShortLinkRepository - FindByCode - Short link found successfully", shared_domain_context.NewField("id", id.String()))
	return pkg.Some(result), nil
}
