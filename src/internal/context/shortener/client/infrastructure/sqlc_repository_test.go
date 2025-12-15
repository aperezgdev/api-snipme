package infrastructure

import (
	"context"
	"testing"

	"github.com/aperezgdev/api-snipme/db/generated"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/internal/context/shortener/client/domain"
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

func TestSqlcClientRepository_Save(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcClientRepository(logger, queries)

	t.Run("Successfully saves a client", func(t *testing.T) {
		client, _ := domain.NewClient("John Doe", "john@example.com")

		err := repo.Save(context.Background(), *client)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err := repo.FindById(context.Background(), client.Id)
		if err != nil {
			t.Fatalf("Expected no error when finding client, got %v", err)
		}
		if !found.IsPresent() {
			t.Fatal("Expected client to be found")
		}
		foundClient := found.Get()
		if foundClient.Id != client.Id {
			t.Errorf("Expected ID %s, got %s", client.Id, foundClient.Id)
		}
		if string(foundClient.Name) != string(client.Name) {
			t.Errorf("Expected Name %s, got %s", client.Name, foundClient.Name)
		}
		if string(foundClient.Email) != string(client.Email) {
			t.Errorf("Expected Email %s, got %s", client.Email, foundClient.Email)
		}
	})

	t.Run("Returns error on duplicate client", func(t *testing.T) {
		client, _ := domain.NewClient("Jane Doe", "jane@example.com")

		err := repo.Save(context.Background(), *client)
		if err != nil {
			t.Fatalf("Expected no error on first save, got %v", err)
		}

		err = repo.Save(context.Background(), *client)
		if err == nil {
			t.Fatal("Expected error on duplicate save, got nil")
		}
	})
}

func TestSqlcClientRepository_FindById(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcClientRepository(logger, queries)

	t.Run("Returns client when found", func(t *testing.T) {
		client, _ := domain.NewClient("Alice", "alice@example.com")
		err := repo.Save(context.Background(), *client)
		if err != nil {
			t.Fatalf("Failed to save client: %v", err)
		}

		found, err := repo.FindById(context.Background(), client.Id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !found.IsPresent() {
			t.Fatal("Expected client to be found")
		}
		foundClient := found.Get()
		if foundClient.Id != client.Id {
			t.Errorf("Expected ID %s, got %s", client.Id, foundClient.Id)
		}
	})

	t.Run("Returns empty when not found", func(t *testing.T) {
		nonExistentId, _ := shared_domain_context.NewID()

		found, err := repo.FindById(context.Background(), nonExistentId)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.IsPresent() {
			t.Error("Expected client not to be found")
		}
	})
}

func TestSqlcClientRepository_Remove(t *testing.T) {
	queries, teardown := setupPostgresContainer(t)
	defer teardown()

	logger := shared_domain_context.DummyLogger{}
	repo := NewSqlcClientRepository(logger, queries)

	t.Run("Successfully removes a client", func(t *testing.T) {
		client, _ := domain.NewClient("Bob", "bob@example.com")
		err := repo.Save(context.Background(), *client)
		if err != nil {
			t.Fatalf("Failed to save client: %v", err)
		}

		found, _ := repo.FindById(context.Background(), client.Id)
		if !found.IsPresent() {
			t.Fatal("Expected client to exist before removal")
		}

		err = repo.Remove(context.Background(), client.Id)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err = repo.FindById(context.Background(), client.Id)
		if err != nil {
			t.Fatalf("Expected no error when finding removed client, got %v", err)
		}
		if found.IsPresent() {
			t.Error("Expected client to be removed")
		}
	})

	t.Run("Does not error when removing non-existent client", func(t *testing.T) {
		nonExistentId, _ := shared_domain_context.NewID()

		err := repo.Remove(context.Background(), nonExistentId)
		if err != nil {
			t.Fatalf("Expected no error when removing non-existent client, got %v", err)
		}
	})
}
