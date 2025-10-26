package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redis.Client
	logger interface {
		Info(ctx context.Context, msg string, fields ...any)
		Error(ctx context.Context, msg string, fields ...any)
	}
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Nuevo constructor con logger
func NewRedisCacheWithLogger(client *redis.Client, logger interface {
	Info(ctx context.Context, msg string, fields ...any)
	Error(ctx context.Context, msg string, fields ...any)
}) *RedisCache {
	return &RedisCache{client: client, logger: logger}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Get - Entrada", "key", key)
	}
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		if c.logger != nil {
			c.logger.Info(ctx, "RedisCache - Get - Clave no encontrada", "key", key)
		}
		return "", nil
	}
	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Get - Error", "key", key, "error", err)
		}
		return "", err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Get - Éxito", "key", key)
	}
	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Set - Entrada", "key", key, "ttl", ttl)
	}
	b, err := json.Marshal(value)
	if err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Set - Error serializando valor", "key", key, "error", err)
		}
		return err
	}
	if err := c.client.Set(ctx, key, b, ttl).Err(); err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Set - Error guardando en Redis", "key", key, "error", err)
		}
		return err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Set - Éxito", "key", key)
	}
	return nil
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Del - Entrada", "key", key)
	}
	if err := c.client.Del(ctx, key).Err(); err != nil {
		if c.logger != nil {
			c.logger.Error(ctx, "RedisCache - Del - Error eliminando clave", "key", key, "error", err)
		}
		return err
	}
	if c.logger != nil {
		c.logger.Info(ctx, "RedisCache - Del - Éxito", "key", key)
	}
	return nil
}
