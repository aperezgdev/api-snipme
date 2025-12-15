package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	link_analytics_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_analytics/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainerForLinkAnalytics(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx := context.Background()

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

CREATE TABLE IF NOT EXISTS link_analytics (
  id uuid PRIMARY KEY,
  link_id uuid NOT NULL UNIQUE,
  total_views INTEGER DEFAULT 0,
  unique_visitors INTEGER DEFAULT 0,
  created_on TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (link_id) REFERENCES short_link(id) ON DELETE CASCADE
);
`)
	if err != nil {
		t.Fatalf("failed to create tables: %s", err)
	}

	cleanup := func() {
		pool.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, cleanup
}

func TestSqlcLinkAnalyticsRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves link analytics", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkAnalytics(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkAnalyticsRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(), 
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "testcode")
		assert.NoError(t, err)

		analyticsId, _ := shared_domain_context.NewID()
		analytics := link_analytics_domain.LinkAnalytics{
			Id:          analyticsId,
			LinkId:      linkId,
			TotalViews:  shared_domain.LinkViewsCounter(10),
			UniqueViews: shared_domain.LinkViewsCounter(5),
			UpdateOn:    shared_domain_context.UpdatedOn(time.Now()),
		}

		err = repo.Save(context.Background(), analytics)
		assert.NoError(t, err)

		found, err := repo.FindByLinkId(context.Background(), linkId)
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
		assert.Equal(t, analyticsId, found.Get().Id)
		assert.Equal(t, linkId, found.Get().LinkId)
		assert.Equal(t, shared_domain.LinkViewsCounter(10), found.Get().TotalViews)
		assert.Equal(t, shared_domain.LinkViewsCounter(5), found.Get().UniqueViews)
	})
}

func TestSqlcLinkAnalyticsRepository_FindByLinkId(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds link analytics by link ID", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkAnalytics(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkAnalyticsRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "findcode")
		assert.NoError(t, err)

		analyticsId, _ := shared_domain_context.NewID()
		analytics := link_analytics_domain.LinkAnalytics{
			Id:          analyticsId,
			LinkId:      linkId,
			TotalViews:  shared_domain.LinkViewsCounter(20),
			UniqueViews: shared_domain.LinkViewsCounter(15),
			UpdateOn:    shared_domain_context.UpdatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), analytics)
		assert.NoError(t, err)

		found, err := repo.FindByLinkId(context.Background(), linkId)
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
		assert.Equal(t, linkId, found.Get().LinkId)
		assert.Equal(t, shared_domain.LinkViewsCounter(20), found.Get().TotalViews)
	})

	t.Run("Returns empty optional when link analytics not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkAnalytics(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkAnalyticsRepository(logger, queries)

		nonExistentId, _ := shared_domain_context.NewID()
		found, err := repo.FindByLinkId(context.Background(), nonExistentId)
		assert.NoError(t, err)
		assert.False(t, found.IsPresent())
	})
}

func TestSqlcLinkAnalyticsRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("Successfully updates link analytics", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkAnalytics(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkAnalyticsRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "updatecode")
		assert.NoError(t, err)

		analyticsId, _ := shared_domain_context.NewID()
		analytics := link_analytics_domain.LinkAnalytics{
			Id:          analyticsId,
			LinkId:      linkId,
			TotalViews:  shared_domain.LinkViewsCounter(10),
			UniqueViews: shared_domain.LinkViewsCounter(5),
			UpdateOn:    shared_domain_context.UpdatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), analytics)
		assert.NoError(t, err)

		analytics.TotalViews = shared_domain.LinkViewsCounter(25)
		analytics.UniqueViews = shared_domain.LinkViewsCounter(18)
		err = repo.Update(context.Background(), analytics)
		assert.NoError(t, err)

		found, err := repo.FindByLinkId(context.Background(), linkId)
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
		assert.Equal(t, shared_domain.LinkViewsCounter(25), found.Get().TotalViews)
		assert.Equal(t, shared_domain.LinkViewsCounter(18), found.Get().UniqueViews)
	})
}
