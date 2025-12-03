package cache

import (
	"context"
	"encoding/json"
	"time"

	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redis.Client
	logger  shared_domain_context.Logger
}

func NewRedisCache(logger shared_domain_context.Logger, client *redis.Client) *RedisCache {
	return &RedisCache{client: client, logger: logger}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Get - ParamsInto", shared_domain_context.NewField("key", key))
	}
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		if c.logger != nil {
			c.logger.Info(ctx, "RedisCache - Get - Key not found", shared_domain_context.NewField("key", key))
		}
		return "", nil
	}
	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Get - Error getting key", shared_domain_context.NewField("key", key), shared_domain_context.NewField("error", err.Error()))
		}
		return "", err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Get - Success", shared_domain_context.NewField("key", key))
	}
	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Set - ParamsInto", shared_domain_context.NewField("key", key), shared_domain_context.NewField("ttl", ttl.String()))
	}
	b, err := json.Marshal(value)
	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Set - Error marshaling value", shared_domain_context.NewField("key", key), shared_domain_context.NewField("error", err.Error()))
		}
		return err
	}
	if err := c.client.Set(ctx, key, b, ttl).Err(); err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Set - Error setting key", shared_domain_context.NewField("key", key), shared_domain_context.NewField("error", err.Error()))
		}
		return err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Set - Success", shared_domain_context.NewField("key", key))
	}
	return nil
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Del - ParamsInto", shared_domain_context.NewField("key", key))
	}
	if err := c.client.Del(ctx, key).Err(); err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Del - Error deleting key", shared_domain_context.NewField("key", key), shared_domain_context.NewField("error", err.Error()))
		}
		return err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Del - Success", shared_domain_context.NewField("key", key))
	}
	return nil
}
