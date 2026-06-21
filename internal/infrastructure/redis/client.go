package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Database struct {
	client *redis.Client
	log    zerolog.Logger
}

func New(i do.Injector) (*Database, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	cfg := do.MustInvoke[*config.Config](i).Redis

	log := logger.WithLayer(rawLog, logger.LayerRedis)
	log.Info().Msg("connecting to redis")

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Error().Err(err).Msg("failed to ping redis")
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	log.Info().Msg("successfully connected to redis")

	return &Database{
		client: client,
		log:    log,
	}, nil
}

func (db *Database) Client() *redis.Client {
	return db.client
}

func (db *Database) HealthCheckWithContext(ctx context.Context) error {
	if err := db.client.Ping(ctx).Err(); err != nil {
		db.log.Error().Err(err).Msg("database health check failed")
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *Database) Shutdown(ctx context.Context) error {
	if db.client != nil {
		db.log.Info().Msg("Closing redis connection")
		_ = db.client.Close()
	}

	return nil
}
