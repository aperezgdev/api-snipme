package infrastructure

import (
	"context"
	"testing"

	"github.com/aperezgdev/api-snipme/db/generated"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupPostgresContainer(t *testing.T) (*generated.Queries, func()) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %s", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %s", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %s", err)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %s", err)
	}

	schema := `
		CREATE TABLE IF NOT EXISTS "user" (
			id UUID PRIMARY KEY,
			oauth_provider VARCHAR(50) NOT NULL,
			oauth_subject VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS client (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			user_id UUID UNIQUE,
			created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT fk_client_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE SET NULL
		);

		CREATE INDEX idx_client_user_id ON client(user_id);

		CREATE TABLE IF NOT EXISTS short_link (
			id uuid PRIMARY KEY,
			summary VARCHAR(255),
			original_route TEXT NOT NULL,
			client_id uuid NULL,
			code VARCHAR(10) UNIQUE NOT NULL,
			created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE
		);
	`
	if _, err := pool.Exec(ctx, schema); err != nil {
		t.Fatalf("Failed to create schema: %s", err)
	}

	queries := generated.New(pool)

	teardown := func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %s", err)
		}
	}

	return queries, teardown
}

func TestSqlcShortLinkRepository_Save(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcShortLinkRepository(logger, queries)

	t.Run("Successfully saves a short link", func(t *testing.T) {
		shortLink, _ := domain.NewShortLink("Example summary", "https://example.com", "")

		err := repo.Save(context.Background(), shortLink)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err := repo.FindById(context.Background(), shortLink.Id)
		if err != nil {
			t.Fatalf("Expected no error when finding short link, got %v", err)
		}
		if !found.IsPresent() {
			t.Fatal("Expected short link to be found")
		}
		foundLink := found.Get()
		if foundLink.Id != shortLink.Id {
			t.Errorf("Expected ID %s, got %s", shortLink.Id, foundLink.Id)
		}
		if string(foundLink.Code) != string(shortLink.Code) {
			t.Errorf("Expected Code %s, got %s", shortLink.Code, foundLink.Code)
		}
		if string(foundLink.OriginalRoute) != string(shortLink.OriginalRoute) {
			t.Errorf("Expected OriginalRoute %s, got %s", shortLink.OriginalRoute, foundLink.OriginalRoute)
		}
	})

	t.Run("Returns error on duplicate code", func(t *testing.T) {
		shortLink1, _ := domain.NewShortLink("First link", "https://example.com/1", "")
		err := repo.Save(context.Background(), shortLink1)
		if err != nil {
			t.Fatalf("Expected no error on first save, got %v", err)
		}

		shortLink2, _ := domain.NewShortLink("Second link", "https://example.com/2", "")
		shortLink2.Code = shortLink1.Code

		err = repo.Save(context.Background(), shortLink2)
		if err == nil {
			t.Fatal("Expected error on duplicate code, got nil")
		}
	})
}

func TestSqlcShortLinkRepository_FindById(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcShortLinkRepository(logger, queries)

	t.Run("Returns short link when found", func(t *testing.T) {
		shortLink, _ := domain.NewShortLink("Example", "https://example.com", "")
		err := repo.Save(context.Background(), shortLink)
		if err != nil {
			t.Fatalf("Failed to save short link: %v", err)
		}

		found, err := repo.FindById(context.Background(), shortLink.Id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !found.IsPresent() {
			t.Fatal("Expected short link to be found")
		}
		foundLink := found.Get()
		if foundLink.Id != shortLink.Id {
			t.Errorf("Expected ID %s, got %s", shortLink.Id, foundLink.Id)
		}
	})

	t.Run("Returns empty when not found", func(t *testing.T) {
		nonExistentId, _ := shared_domain_context.NewID()

		found, err := repo.FindById(context.Background(), nonExistentId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.IsPresent() {
			t.Error("Expected short link not to be found")
		}
	})
}

func TestSqlcShortLinkRepository_FindByCode(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcShortLinkRepository(logger, queries)

	t.Run("Returns short link when found by code", func(t *testing.T) {
		shortLink, _ := domain.NewShortLink("Example", "https://example.com", "")
		err := repo.Save(context.Background(), shortLink)
		if err != nil {
			t.Fatalf("Failed to save short link: %v", err)
		}

		found, err := repo.FindByCode(context.Background(), shortLink.Code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !found.IsPresent() {
			t.Fatal("Expected short link to be found")
		}
		foundLink := found.Get()
		if string(foundLink.Code) != string(shortLink.Code) {
			t.Errorf("Expected Code %s, got %s", shortLink.Code, foundLink.Code)
		}
	})

	t.Run("Returns empty when code not found", func(t *testing.T) {
		code, _ := domain.NewCode()

		found, err := repo.FindByCode(context.Background(), code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.IsPresent() {
			t.Error("Expected short link not to be found")
		}
	})
}

func TestSqlcShortLinkRepository_Remove(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcShortLinkRepository(logger, queries)

	t.Run("Successfully removes a short link", func(t *testing.T) {
		shortLink, _ := domain.NewShortLink("Example", "https://example.com", "")
		err := repo.Save(context.Background(), shortLink)
		if err != nil {
			t.Fatalf("Failed to save short link: %v", err)
		}

		found, _ := repo.FindById(context.Background(), shortLink.Id)
		if !found.IsPresent() {
			t.Fatal("Expected short link to exist before removal")
		}

		err = repo.Remove(context.Background(), shortLink.Id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err = repo.FindById(context.Background(), shortLink.Id)
		if err != nil {
			t.Fatalf("Expected no error when finding removed short link, got %v", err)
		}
		if found.IsPresent() {
			t.Error("Expected short link to be removed")
		}
	})

	t.Run("Does not error when removing non-existent short link", func(t *testing.T) {
		nonExistentId, _ := shared_domain_context.NewID()

		err := repo.Remove(context.Background(), nonExistentId)
		if err != nil {
			t.Fatalf("Expected no error when removing non-existent short link, got %v", err)
		}
	})
}

func TestSqlcShortLinkRepository_FindByClient(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcShortLinkRepository(logger, queries)

	t.Run("Returns short links for a client", func(t *testing.T) {
		clientId, _ := shared_domain_context.NewID()

		// Insert client first to satisfy foreign key constraint
		clientQueries := queries
		clientUUID := pgtype.UUID{}
		_ = clientUUID.Scan(clientId.String())
		createdOn := pgtype.Timestamptz{}
		createdOn.Valid = true

		err := clientQueries.SaveClient(context.Background(), generated.SaveClientParams{
			ID:        clientUUID,
			Name:      "Test Client",
			Email:     "test@example.com",
			CreatedOn: createdOn,
		})
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		shortLink1, _ := domain.NewShortLink("", "https://example.com/1", clientId.String())
		shortLink2, _ := domain.NewShortLink("", "https://example.com/2", clientId.String())

		err = repo.Save(context.Background(), shortLink1)
		if err != nil {
			t.Fatalf("Failed to save first short link: %v", err)
		}
		err = repo.Save(context.Background(), shortLink2)
		if err != nil {
			t.Fatalf("Failed to save second short link: %v", err)
		}

		links, err := repo.FindByClient(context.Background(), clientId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(links) != 2 {
			t.Errorf("Expected 2 links, got %d", len(links))
		}
	})

	t.Run("Returns empty list when client has no links", func(t *testing.T) {
		nonExistentClientId, _ := shared_domain_context.NewID()

		links, err := repo.FindByClient(context.Background(), nonExistentClientId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(links) != 0 {
			t.Errorf("Expected 0 links, got %d", len(links))
		}
	})
}
