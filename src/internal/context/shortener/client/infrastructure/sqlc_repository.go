package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcClientRepository struct {
	queries *generated.Queries
}

func NewSqlcClientRepository(q *generated.Queries) *SqlcClientRepository {
	return &SqlcClientRepository{queries: q}
}

func (r *SqlcClientRepository) Save(ctx context.Context, client domain.Client) error {
	id := pgtype.UUID{}
	_ = id.Scan(client.Id.String())

	createdOn := pgtype.Timestamptz{}
	createdOn.Time = time.Time(client.CreatedOn)
	createdOn.Valid = true

	params := generated.SaveClientParams{
		ID:        id,
		Name:      string(client.Name),
		Email:     string(client.Email),
		CreatedOn: createdOn,
	}
	return r.queries.SaveClient(ctx, params)
}

func (r *SqlcClientRepository) Remove(ctx context.Context, id shared_domain_context.Id) error {
	uuid := pgtype.UUID{}
	_ = uuid.Scan(id.String())
	return r.queries.RemoveClient(ctx, uuid)
}
