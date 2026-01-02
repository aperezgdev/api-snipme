package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/refresh_token/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcRefreshTokenRepository struct {
	logger  shared_domain.Logger
	queries *generated.Queries
}

func NewSqlcRefreshTokenRepository(logger shared_domain.Logger, queries *generated.Queries) *SqlcRefreshTokenRepository {
	return &SqlcRefreshTokenRepository{
		logger:  logger,
		queries: queries,
	}
}

func (r *SqlcRefreshTokenRepository) Save(ctx context.Context, token *domain.RefreshToken) error {
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - Save - Params into",
		shared_domain.NewField("tokenId", token.Id.String()),
		shared_domain.NewField("userId", token.UserId.String()),
		shared_domain.NewField("expiresAt", time.Time(token.ExpiresAt).String()),
	)
	_, err := r.queries.CreateRefreshToken(ctx, generated.CreateRefreshTokenParams{
		ID:        pgtype.UUID{Bytes: token.Id, Valid: true},
		UserID:    pgtype.UUID{Bytes: token.UserId, Valid: true},
		Token:     string(token.Token),
		ExpiresAt: pgtype.Timestamptz{Time: time.Time(token.ExpiresAt), Valid: true},
		CreatedOn: pgtype.Timestamptz{Time: time.Time(token.CreatedOn), Valid: true},
	})

	if err != nil {
		r.logger.Error(ctx, "SqlcRefreshTokenRepository - Save - Error saving refresh token", shared_domain.NewField("error", err.Error()))
		return err
	}

	r.logger.Info(ctx, "SqlcRefreshTokenRepository - Save - Refresh token saved successfully", shared_domain.NewField("tokenId", token.Id.String()))
	return nil
}

func (r *SqlcRefreshTokenRepository) FindByToken(ctx context.Context, token string) (pkg.Optional[*domain.RefreshToken], error) {
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - FindByToken - Searching for refresh token", shared_domain.NewField("token", token))
	refreshToken, err := r.queries.GetRefreshTokenByToken(ctx, token)
	if err != nil {
		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			r.logger.Info(ctx, "SqlcRefreshTokenRepository - FindByToken - Refresh token not found", shared_domain.NewField("token", token))
			return pkg.EmptyOptional[*domain.RefreshToken](), nil
		}
		r.logger.Error(ctx, "SqlcRefreshTokenRepository - FindByToken - Error finding refresh token", shared_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[*domain.RefreshToken](), err
	}

	domainToken := r.toDomain(refreshToken)
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - FindByToken - Refresh token found", shared_domain.NewField("tokenId", domainToken.Id.String()))

	return pkg.Some(domainToken), nil
}

func (r *SqlcRefreshTokenRepository) FindByUserId(ctx context.Context, userId shared_domain.Id) ([]*domain.RefreshToken, error) {
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - FindByUserId - Searching for refresh tokens", shared_domain.NewField("userId", userId.String()))
	tokens, err := r.queries.GetRefreshTokensByUserId(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		r.logger.Error(ctx, "SqlcRefreshTokenRepository - FindByUserId - Error finding refresh tokens", shared_domain.NewField("error", err.Error()))
		return nil, err
	}

	r.logger.Info(ctx, "SqlcRefreshTokenRepository - FindByUserId - Refresh tokens found", shared_domain.NewField("count", len(tokens)))
	domainTokens := pkg.Map(tokens, func(t generated.RefreshToken) *domain.RefreshToken {
		return r.toDomain(t)
	})

	return domainTokens, nil
}

func (r *SqlcRefreshTokenRepository) Delete(ctx context.Context, token string) error {
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - Delete - Deleting refresh token", shared_domain.NewField("token", token))
	err := r.queries.DeleteRefreshToken(ctx, token)
	if err != nil {
		r.logger.Error(ctx, "SqlcRefreshTokenRepository - Delete - Error deleting refresh token", shared_domain.NewField("error", err.Error()))
		return err
	}

	r.logger.Info(ctx, "SqlcRefreshTokenRepository - Delete - Refresh token deleted successfully", shared_domain.NewField("token", token))
	return nil
}

func (r *SqlcRefreshTokenRepository) DeleteByUserId(ctx context.Context, userId shared_domain.Id) error {
	r.logger.Info(ctx, "SqlcRefreshTokenRepository - DeleteByUserId - Deleting refresh tokens by user ID", shared_domain.NewField("userId", userId.String()))
	err := r.queries.DeleteRefreshTokensByUserId(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		r.logger.Error(ctx, "SqlcRefreshTokenRepository - DeleteByUserId - Error deleting refresh tokens", shared_domain.NewField("error", err.Error()))
		return err
	}

	r.logger.Info(ctx, "SqlcRefreshTokenRepository - DeleteByUserId - Refresh tokens deleted successfully", shared_domain.NewField("userId", userId.String()))
	return nil
}

func (r *SqlcRefreshTokenRepository) toDomain(token generated.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		Id:        token.ID.Bytes,
		UserId:    token.UserID.Bytes,
		Token:     domain.RefreshTokenToken(token.Token),
		ExpiresAt: domain.RefreshTokenExpiresAt(token.ExpiresAt.Time),
		CreatedOn: shared_domain.CreatedOn(token.CreatedOn.Time),
	}
}
