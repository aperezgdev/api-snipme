package infrastructure

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/aperezgdev/api-snipme/db/generated"
	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainerForLinkVisit(t *testing.T) (*pgxpool.Pool, func()) {
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

CREATE TABLE IF NOT EXISTS link_visit (
  id uuid PRIMARY KEY,
  link_id uuid NOT NULL,
  visitor_id uuid NOT NULL,
  ip INET,
  user_agent TEXT,
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

func TestSqlcLinkVisitRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("Successfully saves a link visit", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkVisit(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkVisitRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "testcode")
		assert.NoError(t, err)

		visitId, _ := shared_domain_context.NewID()
		visitorId, _ := shared_domain_context.NewID()
		ip := netip.MustParseAddr("192.168.1.1")

		linkVisit := domain.LinkVisit{
			Id:        visitId,
			LinkId:    linkId,
			VisitorId: visitorId,
			Ip:        shared_domain.Ip(ip),
			UserAgent: domain.LinkVisitUserAgent("Mozilla/5.0"),
			CreatedOn: shared_domain_context.CreatedOn(time.Now()),
		}

		err = repo.Save(context.Background(), linkVisit)
		assert.NoError(t, err)

		var count int
		err = pool.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM link_visit WHERE id = $1",
			visitId.String()).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestSqlcLinkVisitRepository_FindOlds(t *testing.T) {
	t.Parallel()

	t.Run("Successfully finds old link visits", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkVisit(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkVisitRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "oldcode")
		assert.NoError(t, err)

		visitId1, _ := shared_domain_context.NewID()
		visitorId1, _ := shared_domain_context.NewID()
		oldDate := time.Now().Add(-20 * time.Minute) // 20 minutes ago

		_, err = pool.Exec(context.Background(),
			`INSERT INTO link_visit (id, link_id, visitor_id, ip, user_agent, created_on) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			visitId1.String(), linkId.String(), visitorId1.String(),
			"192.168.1.1", "Mozilla/5.0", oldDate)
		assert.NoError(t, err)

		visitId2, _ := shared_domain_context.NewID()
		visitorId2, _ := shared_domain_context.NewID()
		recentDate := time.Now().Add(-5 * time.Minute) // 5 minutes ago

		_, err = pool.Exec(context.Background(),
			`INSERT INTO link_visit (id, link_id, visitor_id, ip, user_agent, created_on) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			visitId2.String(), linkId.String(), visitorId2.String(),
			"192.168.1.2", "Chrome/5.0", recentDate)
		assert.NoError(t, err)

		oldVisits, err := repo.FindOlds(context.Background())
		assert.NoError(t, err)
		assert.Len(t, oldVisits, 1)
		assert.Equal(t, visitId1, oldVisits[0].Id)
	})

	t.Run("Returns empty array when no old visits found", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkVisit(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkVisitRepository(logger, queries)

		oldVisits, err := repo.FindOlds(context.Background())
		assert.NoError(t, err)
		assert.Len(t, oldVisits, 0)
	})
}

func TestSqlcLinkVisitRepository_RemoveAll(t *testing.T) {
	t.Parallel()

	t.Run("Successfully removes multiple link visits", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkVisit(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkVisitRepository(logger, queries)

		linkId, _ := shared_domain_context.NewID()
		_, err := pool.Exec(context.Background(),
			"INSERT INTO short_link (id, original_route, code) VALUES ($1, $2, $3)",
			linkId.String(), "https://example.com", "removecode")
		assert.NoError(t, err)

		visitId1, _ := shared_domain_context.NewID()
		visitorId1, _ := shared_domain_context.NewID()
		ip1 := netip.MustParseAddr("192.168.1.1")

		linkVisit1 := domain.LinkVisit{
			Id:        visitId1,
			LinkId:    linkId,
			VisitorId: visitorId1,
			Ip:        shared_domain.Ip(ip1),
			UserAgent: domain.LinkVisitUserAgent("Mozilla/5.0"),
			CreatedOn: shared_domain_context.CreatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), linkVisit1)
		assert.NoError(t, err)

		visitId2, _ := shared_domain_context.NewID()
		visitorId2, _ := shared_domain_context.NewID()
		ip2 := netip.MustParseAddr("192.168.1.2")

		linkVisit2 := domain.LinkVisit{
			Id:        visitId2,
			LinkId:    linkId,
			VisitorId: visitorId2,
			Ip:        shared_domain.Ip(ip2),
			UserAgent: domain.LinkVisitUserAgent("Chrome/5.0"),
			CreatedOn: shared_domain_context.CreatedOn(time.Now()),
		}
		err = repo.Save(context.Background(), linkVisit2)
		assert.NoError(t, err)

		err = repo.RemoveAll(context.Background(), []shared_domain_context.Id{visitId1, visitId2})
		assert.NoError(t, err)

		var count int
		err = pool.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM link_visit WHERE id IN ($1, $2)",
			visitId1.String(), visitId2.String()).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("No error when removing empty list", func(t *testing.T) {
		t.Parallel()

		pool, cleanup := setupPostgresContainerForLinkVisit(t)
		defer cleanup()

		queries := generated.New(pool)
		logger := shared_domain_context.DummyLogger{}
		repo := NewSqlcLinkVisitRepository(logger, queries)

		err := repo.RemoveAll(context.Background(), []shared_domain_context.Id{})
		assert.NoError(t, err)
	})
}
