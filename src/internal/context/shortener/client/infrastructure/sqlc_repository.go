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
	logger  shared_domain_context.Logger
}

func NewSqlcClientRepository(q *generated.Queries, logger shared_domain_context.Logger) *SqlcClientRepository {
	return &SqlcClientRepository{queries: q, logger: logger}
}

func (r *SqlcClientRepository) Save(ctx context.Context, client domain.Client) error {
	r.logger.Info(ctx, "SqlcClientRepository - Save - Params into",
		shared_domain_context.NewField("id", client.Id.String()),
		shared_domain_context.NewField("email", string(client.Email)),
	)
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
	err := r.queries.SaveClient(ctx, params)
	if err != nil {
		r.logger.Error(ctx, "SqlcClientRepository - Save - Error saving client", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcClientRepository - Save - Client saved successfully", shared_domain_context.NewField("id", client.Id.String()))
	return nil
}

func (r *SqlcClientRepository) Remove(ctx context.Context, id shared_domain_context.Id) error {
	r.logger.Info(ctx, "SqlcClientRepository - Remove - Params into", shared_domain_context.NewField("id", id.String()))
	uuid := pgtype.UUID{}
	_ = uuid.Scan(id.String())
	err := r.queries.RemoveClient(ctx, uuid)
	if err != nil {
		r.logger.Error(ctx, "SqlcClientRepository - Remove - Error removing client", shared_domain_context.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcClientRepository - Remove - Client removed successfully", shared_domain_context.NewField("id", id.String()))
	return nil
}
