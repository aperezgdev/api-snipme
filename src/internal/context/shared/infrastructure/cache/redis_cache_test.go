package cache

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
)

func setupRedisContainer(t *testing.T) (container testcontainers.Container, client *goredis.Client, teardown func()) {
	t.Helper()
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get redis host: %v", err)
	}
	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("failed to get redis port: %v", err)
	}
	addr := host + ":" + port.Port()
	client = goredis.NewClient(&goredis.Options{
		Addr: addr,
		DB:   0,
	})
	teardown = func() {
		client.Close()
		container.Terminate(ctx)
	}
	return container, client, teardown
}

func TestRedisCache_Integration(t *testing.T) {
	_, client, teardown := setupRedisContainer(t)
	defer teardown()
	cache := NewRedisCache(&shared_domain_context.DummyLogger{}, client)
	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		key := "foo"
		value := map[string]string{"bar": "baz"}
		err := cache.Set(ctx, key, value, 10*time.Second)
		assert.NoError(t, err)
		val, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.NotEmpty(t, val)
	})

	t.Run("Get missing key returns empty", func(t *testing.T) {
		val, err := cache.Get(ctx, "missing-key")
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	})

	t.Run("Del removes key", func(t *testing.T) {
		key := "todelete"
		cache.Set(ctx, key, "val", 10*time.Second)
		err := cache.Del(ctx, key)
		assert.NoError(t, err)
		val, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	})
}
