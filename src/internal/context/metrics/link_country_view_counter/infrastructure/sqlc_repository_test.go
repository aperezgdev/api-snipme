package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_country_view_counter/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainerForLinkCountryViewCounter(t *testing.T) (*pgxpool.Pool, func()) {
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

	// Create tables programmatically
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

CREATE TABLE IF NOT EXISTS link_country_view_counter (
  id uuid,
  link_id uuid NOT NULL,
  country_code VARCHAR(2) NOT NULL,
  total_views INTEGER DEFAULT 0,
  unique_visitors INTEGER DEFAULT 0,
  created_on TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (link_id) REFERENCES short_link(id) ON DELETE CASCADE,
  PRIMARY KEY (link_id, country_code)
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

func TestSqlcLinkCountryViewCounterRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves a link country view counter", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		// Create a short link first
		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "testcode")
		assert.NoError(t, err)

		// Create link country view counter
		counterId, _ := shared_domain_context.NewID()
		counter := domain.LinkCountryViewCounter{
			Id:          counterId,
			LinkId:      linkId,
			CountryCode: shared_domain.CountryCode("US"),
			TotalViews:  shared_domain.LinkViewsCounter(10),
			UniqueViews: shared_domain.LinkViewsCounter(5),
			CreatedOn:   shared_domain_context.CreatedOn(time.Now()),
		}

		err = repo.Save(context.Background(), counter)
		assert.NoError(t, err)

		// Verify it was saved
		var count int
		err = pool.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM link_country_view_counter WHERE link_id = $1 AND country_code = $2",
			linkId.String(), "US").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestSqlcLinkCountryViewCounterRepository_Find(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds a link country view counter", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		// Create short link
		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "findcode")
		assert.NoError(t, err)

		// Save counter
		counterId, _ := shared_domain_context.NewID()
		counter := domain.LinkCountryViewCounter{
			Id:          counterId,
			LinkId:      linkId,
			CountryCode: shared_domain.CountryCode("ES"),
			TotalViews:  shared_domain.LinkViewsCounter(20),
			UniqueViews: shared_domain.LinkViewsCounter(10),
			CreatedOn:   shared_domain_context.CreatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), counter)
		assert.NoError(t, err)

		// Find counter
		found, err := repo.Find(context.Background(), linkId, shared_domain.CountryCode("ES"))
		assert.NoError(t, err)
		assert.True(t, found.IsPresent())
		assert.Equal(t, linkId, found.Get().LinkId)
		assert.Equal(t, shared_domain.CountryCode("ES"), found.Get().CountryCode)
		assert.Equal(t, shared_domain.LinkViewsCounter(20), found.Get().TotalViews)
	})

	t.Run("Returns empty when country counter not found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()

		found, err := repo.Find(context.Background(), linkId, shared_domain.CountryCode("FR"))
		assert.NoError(t, err)
		assert.False(t, found.IsPresent())
	})
}

func TestSqlcLinkCountryViewCounterRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("Successfully updates a link country view counter", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		// Create short link
		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "updatecode")
		assert.NoError(t, err)

		// Save counter
		counterId, _ := shared_domain_context.NewID()
		counter := domain.LinkCountryViewCounter{
			Id:          counterId,
			LinkId:      linkId,
			CountryCode: shared_domain.CountryCode("DE"),
			TotalViews:  shared_domain.LinkViewsCounter(5),
			UniqueViews: shared_domain.LinkViewsCounter(3),
			CreatedOn:   shared_domain_context.CreatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), counter)
		assert.NoError(t, err)

		// Update counter
		counter.TotalViews = shared_domain.LinkViewsCounter(15)
		counter.UniqueViews = shared_domain.LinkViewsCounter(8)
		err = repo.Update(context.Background(), counter)
		assert.NoError(t, err)

		// Verify update
		var totalViews, uniqueViews int
		err = pool.QueryRow(context.Background(),
			"SELECT total_views, unique_visitors FROM link_country_view_counter WHERE link_id = $1 AND country_code = $2",
			linkId.String(), "DE").Scan(&totalViews, &uniqueViews)
		assert.NoError(t, err)
		assert.Equal(t, 15, totalViews)
		assert.Equal(t, 8, uniqueViews)
	})
}

func TestSqlcLinkCountryViewCounterRepository_RemoveByLink(t *testing.T) {
	t.Parallel()

	t.Run("Successfully removes all counters for a link", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		// Create short link
		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "removecode")
		assert.NoError(t, err)

		// Save multiple counters for different countries
		for _, countryCode := range []string{"US", "ES", "FR"} {
			counterId, _ := shared_domain_context.NewID()
			counter := domain.LinkCountryViewCounter{
				Id:          counterId,
				LinkId:      linkId,
				CountryCode: shared_domain.CountryCode(countryCode),
				TotalViews:  shared_domain.LinkViewsCounter(10),
				UniqueViews: shared_domain.LinkViewsCounter(5),
				CreatedOn:   shared_domain_context.CreatedOn(time.Now()),
			}
			err = repo.Save(context.Background(), counter)
			assert.NoError(t, err)
		}

		// Remove all counters
		err = repo.RemoveByLink(context.Background(), linkId)
		assert.NoError(t, err)

		// Verify they were removed
		var count int
		err = pool.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM link_country_view_counter WHERE link_id = $1",
			linkId.String()).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("No error when removing from non-existent link", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkCountryViewCounter(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkCountryViewCounterRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()

		err := repo.RemoveByLink(context.Background(), linkId)
		assert.NoError(t, err)
	})
}
