package infrastructure

import (
	"context"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SqlcUserRepository struct {
	logger shared_domain.Logger
	queries *generated.Queries
}

func NewSqlcUserRepository(logger shared_domain.Logger, queries *generated.Queries) *SqlcUserRepository {
	return &SqlcUserRepository{
		logger:  logger,
		queries: queries,
	}
}

func (r *SqlcUserRepository) Save(ctx context.Context, user *domain.User) error {
	r.logger.Info(ctx, "SqlcUserRepository - Save - Params into",
		shared_domain.NewField("userId", user.Id.String()),
		shared_domain.NewField("email", user.Email),
		shared_domain.NewField("oauthProvider", user.OAuthProvider.String()),
		shared_domain.NewField("oauthSubject", user.OAuthSubject),
	)
	_, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		ID:            pgtype.UUID{Bytes: user.Id, Valid: true},
		Email:         string(user.Email),
		OauthProvider: user.OAuthProvider.String(),
		OauthSubject:  string(user.OAuthSubject),
		CreatedOn:     pgtype.Timestamptz{Time: time.Time(user.CreatedOn), Valid: true},
	})

	if err != nil {
		r.logger.Error(ctx, "SqlcUserRepository - Save - Error saving user", shared_domain.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcUserRepository - Save - User saved successfully", shared_domain.NewField("userId", user.Id.String()))

	return nil
}

func (r *SqlcUserRepository) FindById(ctx context.Context, id shared_domain.Id) (pkg.Optional[*domain.User], error) {
	r.logger.Info(ctx, "SqlcUserRepository - FindById - Searching for user", shared_domain.NewField("userId", id.String()))
	user, err := r.queries.GetUserById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Info(ctx, "SqlcUserRepository - FindById - User not found", shared_domain.NewField("userId", id.String()))
			return pkg.EmptyOptional[*domain.User](), nil
		}
		
		r.logger.Error(ctx, "SqlcUserRepository - FindById - Error finding user", shared_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[*domain.User](), err
	}

	domainUser := r.toDomain(user)
	r.logger.Info(ctx, "SqlcUserRepository - FindById - User found", shared_domain.NewField("userId", id.String()))

	return pkg.Some(domainUser), nil
}

func (r *SqlcUserRepository) FindByEmail(ctx context.Context, email string) (pkg.Optional[*domain.User], error) {
	r.logger.Info(ctx, "SqlcUserRepository - FindByEmail - Searching for user", shared_domain.NewField("email", email))
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Info(ctx, "SqlcUserRepository - FindByEmail - User not found", shared_domain.NewField("email", email))
			return pkg.EmptyOptional[*domain.User](), nil
		}
		r.logger.Error(ctx, "SqlcUserRepository - FindByEmail - Error finding user", shared_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[*domain.User](), err
	}

	domainUser := r.toDomain(user)

	r.logger.Info(ctx, "SqlcUserRepository - FindByEmail - User found", shared_domain.NewField("email", email))
	return pkg.Some(domainUser), nil
}

func (r *SqlcUserRepository) FindByOAuthProviderAndSubject(
	ctx context.Context,
	provider domain.OAuthProvider,
	subject string,
) (pkg.Optional[*domain.User], error) {
	r.logger.Info(ctx, "SqlcUserRepository - FindByOAuthProviderAndSubject - Searching for user", shared_domain.NewField("provider", provider.String()), shared_domain.NewField("subject", subject))
	user, err := r.queries.GetUserByOAuthProviderAndSubject(ctx, generated.GetUserByOAuthProviderAndSubjectParams{
		OauthProvider: provider.String(),
		OauthSubject:  subject,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Info(ctx, "SqlcUserRepository - FindByOAuthProviderAndSubject - User not found", shared_domain.NewField("provider", provider.String()), shared_domain.NewField("subject", subject))
			return pkg.EmptyOptional[*domain.User](), nil
		}
		r.logger.Error(ctx, "SqlcUserRepository - FindByOAuthProviderAndSubject - Error finding user", shared_domain.NewField("error", err.Error()))
		return pkg.EmptyOptional[*domain.User](), err
	}

	domainUser := r.toDomain(user)
	r.logger.Info(ctx, "SqlcUserRepository - FindByOAuthProviderAndSubject - User found", shared_domain.NewField("userId", domainUser.Id.String()))

	return pkg.Some(domainUser), nil
}

func (r *SqlcUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.logger.Info(ctx, "SqlcUserRepository - Update - Updating user", shared_domain.NewField("userId", user.Id.String()))
	_, err := r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:            pgtype.UUID{Bytes: user.Id, Valid: true},
		Email:         string(user.Email),
		OauthProvider: user.OAuthProvider.String(),
		OauthSubject:  string(user.OAuthSubject),
	})
	if err != nil {
		r.logger.Error(ctx, "SqlcUserRepository - Update - Error updating user", shared_domain.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "SqlcUserRepository - Update - User updated successfully", shared_domain.NewField("userId", user.Id.String()))

	return err
}

func (r *SqlcUserRepository) toDomain(user generated.User) *domain.User {
	return &domain.User{
		Id:            	user.ID.Bytes,
		Email:          shared_domain.Email(user.Email),
		OAuthProvider:  domain.OAuthProvider(user.OauthProvider),
		OAuthSubject:  	domain.OAuthSubject(user.OauthSubject),
		CreatedOn:			shared_domain.CreatedOn(user.CreatedOn.Time),
	}
}
