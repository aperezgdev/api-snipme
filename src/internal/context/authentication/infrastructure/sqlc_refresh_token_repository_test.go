package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainerForRefreshToken(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts(
			"../../../../../db/schema/user.sql",
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create connection pool: %s", err)
	}

	// Create refresh_token table after container is ready
	_, err = pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS refresh_token (
id UUID PRIMARY KEY,
user_id UUID NOT NULL,
token TEXT NOT NULL UNIQUE,
expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
CONSTRAINT fk_refresh_token_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_refresh_token_user_id ON refresh_token(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_token_token ON refresh_token(token);
CREATE INDEX IF NOT EXISTS idx_refresh_token_expires_at ON refresh_token(expires_at);
`)
	if err != nil {
		t.Fatalf("failed to create refresh_token table: %s", err)
	}

	cleanup := func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}

func TestSqlcRefreshTokenRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves a refresh token", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("token@example.com", domain.OAuthProviderGoogle, "token-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		refreshToken, err := domain.NewRefreshToken(user.Id.String(), "test-token-123", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		err = tokenRepo.Save(context.Background(), refreshToken)
		assert.NoError(t, err)

		foundToken, err := tokenRepo.FindByToken(context.Background(), "test-token-123")
		assert.NoError(t, err)
		assert.True(t, foundToken.IsPresent())
		assert.Equal(t, refreshToken.Id, foundToken.Get().Id)
		assert.Equal(t, user.Id, foundToken.Get().UserId)
	})

	t.Run("Returns error on duplicate token", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("duplicate-token@example.com", domain.OAuthProviderGoogle, "dup-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		token1, err := domain.NewRefreshToken(user.Id.String(), "duplicate-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token1)
		assert.NoError(t, err)

		token2, err := domain.NewRefreshToken(user.Id.String(), "duplicate-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token2)
		assert.Error(t, err)
	})
}
func TestSqlcRefreshTokenRepository_FindByToken(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds a refresh token by token string", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("findtoken@example.com", domain.OAuthProviderGoogle, "find-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		refreshToken, err := domain.NewRefreshToken(user.Id.String(), "findable-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), refreshToken)
		assert.NoError(t, err)

		foundToken, err := tokenRepo.FindByToken(context.Background(), "findable-token")
		assert.NoError(t, err)
		assert.True(t, foundToken.IsPresent())
		assert.Equal(t, refreshToken.Id, foundToken.Get().Id)
		assert.Equal(t, user.Id, foundToken.Get().UserId)
		assert.Equal(t, "findable-token", string(foundToken.Get().Token))
	})

	t.Run("Returns empty optional when token not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		foundToken, err := tokenRepo.FindByToken(context.Background(), "non-existent-token")
		assert.NoError(t, err)
		assert.False(t, foundToken.IsPresent())
	})
}

func TestSqlcRefreshTokenRepository_FindByUserId(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds all tokens for a user", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("multitoken@example.com", domain.OAuthProviderGitHub, "multi-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		token1, err := domain.NewRefreshToken(user.Id.String(), "user-token-1", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token1)
		assert.NoError(t, err)

		token2, err := domain.NewRefreshToken(user.Id.String(), "user-token-2", time.Now().Add(48*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token2)
		assert.NoError(t, err)

		tokens, err := tokenRepo.FindByUserId(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.Len(t, tokens, 2)
	})

	t.Run("Returns empty array when no tokens found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("notoken@example.com", domain.OAuthProviderGitHub, "no-token-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		tokens, err := tokenRepo.FindByUserId(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.Len(t, tokens, 0)
	})
}

func TestSqlcRefreshTokenRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Successfully deletes a refresh token", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("delete@example.com", domain.OAuthProviderGoogle, "delete-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		refreshToken, err := domain.NewRefreshToken(user.Id.String(), "delete-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), refreshToken)
		assert.NoError(t, err)

		err = tokenRepo.Delete(context.Background(), "delete-token")
		assert.NoError(t, err)

		foundToken, err := tokenRepo.FindByToken(context.Background(), "delete-token")
		assert.NoError(t, err)
		assert.False(t, foundToken.IsPresent())
	})

	t.Run("No error when deleting non-existent token", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		err := tokenRepo.Delete(context.Background(), "non-existent-token")
		assert.NoError(t, err)
	})
}

func TestSqlcRefreshTokenRepository_DeleteByUserId(t *testing.T) {
	t.Parallel()

	t.Run("Successfully deletes all tokens for a user", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("deleteall@example.com", domain.OAuthProviderGitHub, "deleteall-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		token1, err := domain.NewRefreshToken(user.Id.String(), "delete-user-token-1", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token1)
		assert.NoError(t, err)

		token2, err := domain.NewRefreshToken(user.Id.String(), "delete-user-token-2", time.Now().Add(48*time.Hour))
		assert.NoError(t, err)
		err = tokenRepo.Save(context.Background(), token2)
		assert.NoError(t, err)

		err = tokenRepo.DeleteByUserId(context.Background(), user.Id)
		assert.NoError(t, err)

		tokens, err := tokenRepo.FindByUserId(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.Len(t, tokens, 0)
	})

	t.Run("No error when deleting tokens for user with no tokens", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForRefreshToken(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		userRepo := NewSqlcUserRepository(logger, queries)
		tokenRepo := NewSqlcRefreshTokenRepository(logger, queries)

		user, err := domain.NewUser("notokensdelete@example.com", domain.OAuthProviderGitHub, "no-tokens-delete-subject")
		assert.NoError(t, err)
		err = userRepo.Save(context.Background(), user)
		assert.NoError(t, err)

		err = tokenRepo.DeleteByUserId(context.Background(), user.Id)
		assert.NoError(t, err)
	})
}
