package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	redisInfra "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/redis"
	"github.com/samber/do/v2"
)

type TokenRepository interface {
	IsBlacklisted(ctx context.Context, accessTokenID uuid.UUID) bool
	Blacklist(ctx context.Context, accessTokenID uuid.UUID, expiredAt time.Time) error
}

type tokenRepository struct {
	client *redis.Client
	conf   *config.Config
}

func NewTokenRepository(i do.Injector) (TokenRepository, error) {
	client := do.MustInvoke[*redisInfra.Database](i)
	conf := do.MustInvoke[*config.Config](i)

	return tokenRepository{
		client: client.Client(),
		conf:   conf,
	}, nil
}

func (r tokenRepository) KeyBlacklist(accessTokenID uuid.UUID) string {
	return fmt.Sprintf("auth:blacklist:%s", accessTokenID)
}

func (r tokenRepository) IsBlacklisted(ctx context.Context, accessTokenID uuid.UUID) bool {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)
	key := r.KeyBlacklist(accessTokenID)
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Error().Err(err).Msg("failed to check if token exists in blacklist")
		}

		return false
	}

	return result > 0
}

func (r tokenRepository) Blacklist(ctx context.Context, accessTokenID uuid.UUID, expiredAt time.Time) error {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)

	key := r.KeyBlacklist(accessTokenID)
	err := r.client.Set(ctx, key, "1", time.Until(expiredAt)).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to add token to blacklist")
		return err
	}

	return nil
}
