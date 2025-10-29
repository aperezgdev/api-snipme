package test

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/aperezgdev/api-snipme/src/cmd/bootstrap"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var API_URL = "http://localhost:8081"

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpassword"),
		postgres.BasicWaitStrategies(),
		testcontainers.WithExposedPorts("5431"),
	)

	containerRedis, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})

	if err != nil {
		panic(err)
	}
	defer containerRedis.Terminate(ctx)

	redisHost, err := containerRedis.Host(ctx)
	if err != nil {
		panic(err)
	}
	redisPort, err := containerRedis.MappedPort(ctx, "6379")
	if err != nil {
		panic(err)
	}
	os.Setenv("REDIS_URL", redisHost+":"+redisPort.Port())
	defer container.Terminate(ctx)

	connStr, _ := container.ConnectionString(ctx, "sslmode=disable")

	os.Setenv("DATABASE_URL", connStr)
	os.Setenv("SERVER_PORT", "8081")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	os.Setenv("ENV", "test")

	schemaFiles := []string{
		"../db/schema/client.sql",
		"../db/schema/short_link.sql",
		"../db/schema/link_analytics.sql",
		"../db/schema/link_country_view_counter.sql",
		"../db/schema/link_visit.sql",
		"../db/seed/test_seed.sql",
	}

	for _, file := range schemaFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			panic(err)
		}
	}

	os.Setenv("DATABASE_URL", connStr)
	os.Setenv("PORT", "8081")

	go func() {
		err := bootstrap.Run()
		if err != nil {
			panic(err)
		}
	}()

	waitForServer(API_URL + "/status")

	code := m.Run()
	os.Exit(code)
}

func waitForServer(url string) {
	for range 3 {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
	panic("API server not responding")
}
