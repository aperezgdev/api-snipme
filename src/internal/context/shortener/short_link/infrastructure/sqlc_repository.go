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
	queries *generated.Queries
}

func NewSqlcShortLinkRepository(q *generated.Queries) *SqlcShortLinkRepository {
	return &SqlcShortLinkRepository{queries: q}
}

func (r *SqlcShortLinkRepository) Save(ctx context.Context, shortLink *domain.ShortLink) error {
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

	return err
}

func (r *SqlcShortLinkRepository) Remove(ctx context.Context, id shared_domain_context.Id) error {
	uuid := pgtype.UUID{}
	_ = uuid.Scan(id.String())
	return r.queries.RemoveShortLink(ctx, uuid)
}

func (r *SqlcShortLinkRepository) FindByCode(ctx context.Context, code domain.ShortLinkCode) (pkg.Optional[*domain.ShortLink], error) {
	shortLink, err := r.queries.FindShortLinkByCode(ctx, string(code))
	if err != nil {
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
	return pkg.Some(result), nil
}
