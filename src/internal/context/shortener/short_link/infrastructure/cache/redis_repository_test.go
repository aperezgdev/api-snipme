package cache

import (
	"context"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_cache "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/cache"
	shortlink_domain "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	shortlink_infra "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/infrastructure"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisAndPostgres(t *testing.T) (*goredis.Client, *pgxpool.Pool, func()) {
	t.Helper()
	ctx := context.Background()

	// Setup Redis
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("failed to get redis port: %v", err)
	}

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: redisHost + ":" + redisPort.Port(),
		DB:   0,
	})

	// Setup PostgreSQL
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
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

	// Create tables
	_, err = pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS client (
  id uuid PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS short_link (
  id uuid PRIMARY KEY,
  summary VARCHAR(255),
  original_route TEXT NOT NULL,
  client_id uuid NULL,
  code VARCHAR(10) UNIQUE NOT NULL,
  created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE
);
`)
	if err != nil {
		t.Fatalf("failed to create tables: %s", err)
	}

	cleanup := func() {
		redisClient.Close()
		pool.Close()
		redisContainer.Terminate(ctx)
		postgresContainer.Terminate(ctx)
	}

	return redisClient, pool, cleanup
}

func TestRedisShortLinkRepository_FindByCode(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds from cache on second call", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		shortLink, _ := shortlink_domain.NewPublicShortLink("https://example.com")
		err := redisRepo.Save(context.Background(), shortLink)
		assert.NoError(t, err)

		found1, err := redisRepo.FindByCode(context.Background(), shortLink.Code)
		assert.NoError(t, err)
		assert.True(t, found1.IsPresent())
		assert.Equal(t, shortLink.Id, found1.Get().Id)

		found2, err := redisRepo.FindByCode(context.Background(), shortLink.Code)
		assert.NoError(t, err)
		assert.True(t, found2.IsPresent())
		assert.Equal(t, shortLink.Id, found2.Get().Id)

		key := "shortlink:" + string(shortLink.Code)
		val, err := redisClient.Get(context.Background(), key).Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, val)
	})

	t.Run("Returns empty optional when not found", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		code := shortlink_domain.ShortLinkCode("nonexist")
		found, err := redisRepo.FindByCode(context.Background(), code)
		assert.NoError(t, err)
		assert.False(t, found.IsPresent())
	})
}

func TestRedisShortLinkRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves and invalidates cache", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		shortLink, _ := shortlink_domain.NewPublicShortLink("https://example.com")

		key := "shortlink:" + string(shortLink.Code)
		redisClient.Set(context.Background(), key, "oldvalue", 5*time.Minute)

		err := redisRepo.Save(context.Background(), shortLink)
		assert.NoError(t, err)

		val, err := redisClient.Get(context.Background(), key).Result()
		assert.Error(t, err)
		assert.Empty(t, val)

		found, err := redisRepo.FindById(context.Background(), shortLink.Id)
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
	})
}

func TestRedisShortLinkRepository_FindById(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds by ID", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		shortLink, _ := shortlink_domain.NewPublicShortLink("https://example.com")
		err := redisRepo.Save(context.Background(), shortLink)
		assert.NoError(t, err)

		found, err := redisRepo.FindById(context.Background(), shortLink.Id)
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
		assert.Equal(t, shortLink.Id, found.Get().Id)
	})

	t.Run("Returns empty optional when not found", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		nonExistentId, _ := domain.NewID()
		found, err := redisRepo.FindById(context.Background(), nonExistentId)
		assert.NoError(t, err)
		assert.False(t, found.IsPresent())
	})
}

func TestRedisShortLinkRepository_Remove(t *testing.T) {
	t.Parallel()

	t.Run("Successfully removes short link", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		shortLink, _ := shortlink_domain.NewPublicShortLink("https://example.com")
		err := redisRepo.Save(context.Background(), shortLink)
		assert.NoError(t, err)

		err = redisRepo.Remove(context.Background(), shortLink.Id)
		assert.NoError(t, err)

		found, err := redisRepo.FindById(context.Background(), shortLink.Id)
		assert.NoError(t, err)
		assert.False(t, found.IsPresent())
	})
}

func TestRedisShortLinkRepository_FindByClient(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds all links by client ID", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		clientId, _ := domain.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO client (id, name, email) VALUES ($1, $2, $3)",
			clientId.String(), "Test Client", "test@client.com")
		assert.NoError(t, err)

		shortLink1, _ := shortlink_domain.NewShortLink("", "https://example1.com", clientId.String())
		err = redisRepo.Save(context.Background(), shortLink1)
		assert.NoError(t, err)

		shortLink2, _ := shortlink_domain.NewShortLink("", "https://example2.com", clientId.String())
		err = redisRepo.Save(context.Background(), shortLink2)
		assert.NoError(t, err)

		links, err := redisRepo.FindByClient(context.Background(), clientId)
		assert.NoError(t, err)
		assert.Len(t, links, 2)
	})

	t.Run("Returns empty array when no links found", func(t *testing.T) {
		t.Parallel()

		redisClient, pool, cleanup := setupRedisAndPostgres(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := domain.DummyLogger{}
		sqlRepo := shortlink_infra.NewSqlcShortLinkRepository(logger, queries)
		cache := shared_cache.NewRedisCache(&logger, redisClient)
		redisRepo := NewRedisShortLinkRepository(sqlRepo, cache, 5*time.Minute, logger)

		nonExistentClientId, _ := domain.NewID()
		links, err := redisRepo.FindByClient(context.Background(), nonExistentClientId)
		assert.NoError(t, err)
		assert.Len(t, links, 0)
	})
}
