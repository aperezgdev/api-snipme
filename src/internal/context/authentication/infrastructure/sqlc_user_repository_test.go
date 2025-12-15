package infrastructure

import (
	"context"
	"testing"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/authentication/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainerForUser(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("../../../../../db/schema/user.sql"),
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

	cleanup := func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}

func TestSqlcUserRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves a user", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user, err := domain.NewUser("test@example.com", domain.OAuthProviderGoogle, "oauth-subject-123")
		assert.NoError(t, err)

		err = repo.Save(context.Background(), user)
		assert.NoError(t, err)

		foundUser, err := repo.FindById(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.True(t, foundUser.IsPresent())
		assert.Equal(t, user.Id, foundUser.Get().Id)
		assert.Equal(t, user.Email, foundUser.Get().Email)
	})

	t.Run("Returns error on duplicate email", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user1, err := domain.NewUser("duplicate@example.com", domain.OAuthProviderGoogle, "subject-1")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user1)
		assert.NoError(t, err)

		user2, err := domain.NewUser("duplicate@example.com", domain.OAuthProviderGoogle, "subject-2")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user2)
		assert.Error(t, err)
	})

	t.Run("Returns error on duplicate oauth identity", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user1, err := domain.NewUser("user1@example.com", domain.OAuthProviderGoogle, "same-subject")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user1)
		assert.NoError(t, err)

		user2, err := domain.NewUser("user2@example.com", domain.OAuthProviderGoogle, "same-subject")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user2)
		assert.Error(t, err)
	})
}

func TestSqlcUserRepository_FindById(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds a user by ID", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user, err := domain.NewUser("findbyid@example.com", domain.OAuthProviderGoogle, "subject-findbyid")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user)
		assert.NoError(t, err)

		foundUser, err := repo.FindById(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.True(t, foundUser.IsPresent())
		assert.Equal(t, user.Id, foundUser.Get().Id)
		assert.Equal(t, user.Email, foundUser.Get().Email)
		assert.Equal(t, user.OAuthProvider, foundUser.Get().OAuthProvider)
		assert.Equal(t, user.OAuthSubject, foundUser.Get().OAuthSubject)
	})

	t.Run("Returns empty optional when user not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		nonExistentId, _ := shared_domain.NewID()
		foundUser, err := repo.FindById(context.Background(), nonExistentId)
		assert.NoError(t, err)
		assert.False(t, foundUser.IsPresent())
	})
}

func TestSqlcUserRepository_FindByEmail(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds a user by email", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user, err := domain.NewUser("findbyemail@example.com", domain.OAuthProviderGitHub, "subject-findbyemail")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user)
		assert.NoError(t, err)

		foundUser, err := repo.FindByEmail(context.Background(), "findbyemail@example.com")
		assert.NoError(t, err)
		assert.True(t, foundUser.IsPresent())
		assert.Equal(t, user.Id, foundUser.Get().Id)
		assert.Equal(t, user.Email, foundUser.Get().Email)
	})

	t.Run("Returns empty optional when user not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		foundUser, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
		assert.NoError(t, err)
		assert.False(t, foundUser.IsPresent())
	})
}

func TestSqlcUserRepository_FindByOAuthProviderAndSubject(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds a user by OAuth provider and subject", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user, err := domain.NewUser("oauth@example.com", domain.OAuthProviderGoogle, "oauth-subject-unique")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user)
		assert.NoError(t, err)

		foundUser, err := repo.FindByOAuthProviderAndSubject(context.Background(), domain.OAuthProviderGoogle, "oauth-subject-unique")
		assert.NoError(t, err)
		assert.True(t, foundUser.IsPresent())
		assert.Equal(t, user.Id, foundUser.Get().Id)
		assert.Equal(t, user.OAuthProvider, foundUser.Get().OAuthProvider)
		assert.Equal(t, user.OAuthSubject, foundUser.Get().OAuthSubject)
	})

	t.Run("Returns empty optional when user not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		foundUser, err := repo.FindByOAuthProviderAndSubject(context.Background(), domain.OAuthProviderGitHub, "nonexistent-subject")
		assert.NoError(t, err)
		assert.False(t, foundUser.IsPresent())
	})
}

func TestSqlcUserRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("Successfully updates a user", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForUser(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain.DummyLogger{}
		repo := NewSqlcUserRepository(logger, queries)

		user, err := domain.NewUser("update@example.com", domain.OAuthProviderGoogle, "update-subject")
		assert.NoError(t, err)
		err = repo.Save(context.Background(), user)
		assert.NoError(t, err)

		newEmail, _ := shared_domain.NewEmail("updated@example.com")
		user.Email = newEmail

		err = repo.Update(context.Background(), user)
		assert.NoError(t, err)

		foundUser, err := repo.FindById(context.Background(), user.Id)
		assert.NoError(t, err)
		assert.True(t, foundUser.IsPresent())
		assert.Equal(t, newEmail, foundUser.Get().Email)
	})
}
