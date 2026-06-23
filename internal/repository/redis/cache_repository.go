package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	redisInfra "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/redis"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
	"github.com/samber/do/v2"
)

type CacheRepository interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type cacheRepository struct {
	client *redis.Client
	conf   *config.Config
}

func NewCacheRepository(i do.Injector) (CacheRepository, error) {
	client := do.MustInvoke[*redisInfra.Database](i)
	conf := do.MustInvoke[*config.Config](i)

	return cacheRepository{
		client: client.Client(),
		conf:   conf,
	}, nil
}

func (r cacheRepository) Get(ctx context.Context, key string, dest any) error {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)

	if err := helper.MustBePointer[any](dest); err != nil {
		return err
	}

	if key == "" {
		return fmt.Errorf("cache key cannot be empty")
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			log.Error().Err(err).Msg("failed to get key from redis")
		}

		return err
	}

	if err := sonic.Unmarshal(val, dest); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal redis data")
		return err
	}

	return nil
}

func (r cacheRepository) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)

	if key == "" {
		return fmt.Errorf("cache key cannot be empty")
	}

	val, err := sonic.Marshal(value)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal data for redis")
		return err
	}

	if err := r.client.Set(ctx, key, val, expiration).Err(); err != nil {
		log.Error().Err(err).Msg("failed to set key in redis")
		return err
	}

	return nil
}

func (r cacheRepository) Delete(ctx context.Context, key string) error {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)

	if key == "" {
		return fmt.Errorf("cache key cannot be empty")
	}

	if err := r.client.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		log.Error().Err(err).Msg("failed to delete key from redis")
		return err
	}

	return nil
}
