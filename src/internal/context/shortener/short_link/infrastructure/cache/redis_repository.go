package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	shared_cache "github.com/aperezgdev/api-snipme/src/internal/context/shared/infrastructure/cache"
	shortlink_domain "github.com/aperezgdev/api-snipme/src/internal/context/shortener/short_link/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
)

type RedisShortLinkRepository struct {
	repo      shortlink_domain.ShortLinkRepository
	cache     shared_cache.Cache
	ttl       time.Duration
	keyPrefix string
	logger    domain.Logger
}

func NewRedisShortLinkRepository(
	repo shortlink_domain.ShortLinkRepository,
	cache shared_cache.Cache,
	ttl time.Duration,
	logger domain.Logger,
) *RedisShortLinkRepository {
	return &RedisShortLinkRepository{
		repo:      repo,
		cache:     cache,
		ttl:       ttl,
		keyPrefix: "shortlink:",
		logger:    logger,
	}
}

func (r *RedisShortLinkRepository) Save(ctx context.Context, shortLink *shortlink_domain.ShortLink) error {
	r.logger.Info(ctx, "RedisShortLinkRepository - Save - Params into", domain.NewField("code", string(shortLink.Code)))
	if err := r.cache.Del(ctx, r.codeKey(shortLink.Code)); err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - Save - Error deleting cache", domain.NewField("error", err.Error()))
		return err
	}
	err := r.repo.Save(ctx, shortLink)
	if err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - Save - Error saving to repo", domain.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "RedisShortLinkRepository - Save - Success", domain.NewField("code", string(shortLink.Code)))
	return nil
}

func (r *RedisShortLinkRepository) FindById(ctx context.Context, id domain.Id) (pkg.Optional[*shortlink_domain.ShortLink], error) {
	r.logger.Info(ctx, "RedisShortLinkRepository - FindById - Params into", domain.NewField("id", id.String()))
	result, err := r.repo.FindById(ctx, id)
	if err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - FindById - Error finding in repo", domain.NewField("error", err.Error()))
		return pkg.Optional[*shortlink_domain.ShortLink]{}, err
	}
	r.logger.Info(ctx, "RedisShortLinkRepository - FindById - Success", domain.NewField("id", id.String()))
	return result, nil
}

func (r *RedisShortLinkRepository) Remove(ctx context.Context, id domain.Id) error {
	r.logger.Info(ctx, "RedisShortLinkRepository - Remove - Params into", domain.NewField("id", id.String()))
	err := r.repo.Remove(ctx, id)
	if err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - Remove - Error removing from repo", domain.NewField("error", err.Error()))
		return err
	}
	r.logger.Info(ctx, "RedisShortLinkRepository - Remove - Success", domain.NewField("id", id.String()))
	return nil
}

func (r *RedisShortLinkRepository) FindByCode(ctx context.Context, code shortlink_domain.ShortLinkCode) (pkg.Optional[*shortlink_domain.ShortLink], error) {
	r.logger.Info(ctx, "RedisShortLinkRepository - FindByCode - Params into", domain.NewField("code", string(code)))
	key := r.codeKey(code)
	val, err := r.cache.Get(ctx, key)
	if err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - FindByCode - Error getting from cache", domain.NewField("error", err.Error()))
		return pkg.Optional[*shortlink_domain.ShortLink]{}, err
	}
	if val != "" {
		var cached shortlink_domain.ShortLink
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			r.logger.Info(ctx, "RedisShortLinkRepository - FindByCode - Cache hit", domain.NewField("code", string(code)))
			return pkg.Some(&cached), nil
		} else {
			r.logger.Error(ctx, "RedisShortLinkRepository - FindByCode - Error unmarshalling cache", domain.NewField("error", err.Error()))
			if errDel := r.cache.Del(ctx, key); errDel != nil {
				r.logger.Error(ctx, "RedisShortLinkRepository - FindByCode - Error deleting corrupted cache", domain.NewField("error", errDel.Error()))
				return pkg.Optional[*shortlink_domain.ShortLink]{}, errDel
			}
			return pkg.Optional[*shortlink_domain.ShortLink]{}, err
		}
	}
	result, err := r.repo.FindByCode(ctx, code)
	if err != nil || !result.IsPresent() {
		if err != nil {
			r.logger.Error(ctx, "RedisShortLinkRepository - FindByCode - Error finding in repo", domain.NewField("error", err.Error()))
		}
		return result, err
	}
	if err := r.cache.Set(ctx, key, result.Get(), r.ttl); err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - FindByCode - Error setting cache", domain.NewField("error", err.Error()))
		return result, err
	}
	r.logger.Info(ctx, "RedisShortLinkRepository - FindByCode - Success", domain.NewField("code", string(code)))
	return result, nil
}

func (r *RedisShortLinkRepository) FindByClient(ctx context.Context, clientId domain.Id) ([]*shortlink_domain.ShortLink, error) {
	r.logger.Info(ctx, "RedisShortLinkRepository - FindByClient - Params into", domain.NewField("clientId", clientId.String()))
	links, err := r.repo.FindByClient(ctx, clientId)
	if err != nil {
		r.logger.Error(ctx, "RedisShortLinkRepository - FindByClient - Error finding in repo", domain.NewField("error", err.Error()))
		return nil, err
	}
	r.logger.Info(ctx, "RedisShortLinkRepository - FindByClient - Success", domain.NewField("count", len(links)))
	return links, nil
}

func (r *RedisShortLinkRepository) codeKey(code shortlink_domain.ShortLinkCode) string {
	return fmt.Sprintf("%s%s", r.keyPrefix, code)
}
