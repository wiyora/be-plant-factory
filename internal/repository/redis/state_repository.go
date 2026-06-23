package redis

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	redisInfra "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/redis"
	"github.com/samber/do/v2"
)

type StateRepository interface {
	Set(ctx context.Context, provider entity.Provider, state StatePayload) error
	ValidateAndDelete(ctx context.Context, provider entity.Provider, state string) (StatePayload, error)
}

type stateRepository struct {
	client *redis.Client
	conf   *config.Config
}

func NewStateRepository(i do.Injector) (StateRepository, error) {
	client := do.MustInvoke[*redisInfra.Database](i)
	conf := do.MustInvoke[*config.Config](i)

	return stateRepository{
		client: client.Client(),
		conf:   conf,
	}, nil
}

func (r stateRepository) Key(provider entity.Provider, state string) string {
	return fmt.Sprintf("auth:%s:%s", provider, state)
}

func (r stateRepository) Set(ctx context.Context, provider entity.Provider, state StatePayload) error {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)
	key := r.Key(provider, state.State)

	val, err := sonic.Marshal(state)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal data for redis")
		return err
	}

	if err := r.client.Set(ctx, key, val, r.conf.Auth.SocialStateDuration).Err(); err != nil {
		log.Error().Err(err).Msg("failed to execute query")
		return err
	}

	return nil
}

func (r stateRepository) ValidateAndDelete(ctx context.Context, provider entity.Provider, state string) (StatePayload, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerRedisRepository)
	key := r.Key(provider, state)

	val, err := r.client.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return StatePayload{}, nil
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to execute query")
		return StatePayload{}, err
	}

	var dest StatePayload
	if err := sonic.Unmarshal([]byte(val), &dest); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal data from redis")
		return StatePayload{}, err
	}

	return dest, nil
}
