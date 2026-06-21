package postgres

import (
	"context"
	"fmt"

	pgxDecimal "github.com/ColeBurch/pgx-govalues-decimal"
	pgxZerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Database struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func New(i do.Injector) (*Database, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	cfg := do.MustInvoke[*config.Config](i)
	conf := cfg.Database

	log := logger.WithLayer(rawLog, logger.LayerPostgres)
	log.Info().Msg("connecting to postgres")

	poolConfig, err := pgxpool.ParseConfig(conf.DSN())
	if err != nil {
		log.Error().Err(err).Msg("failed to parse postgres config")
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(conf.MaxOpenConns)
	poolConfig.MinConns = int32(conf.MinIdleConns)
	poolConfig.MaxConnLifetime = conf.MaxConnLifetime
	poolConfig.MaxConnIdleTime = conf.MaxConnIdleTime

	if cfg.App.Env.IsDevelopment() {
		logger := pgxZerolog.NewLogger(log, pgxZerolog.WithoutPGXModule())

		poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   logger,
			LogLevel: tracelog.LogLevelTrace,
		}
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxDecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create postgres pool")
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		log.Error().Err(err).Msg("failed to ping postgres")
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	log.Info().Msg("successfully connected to postgres")

	if err := AutoMigrate(conf.DSN(), log); err != nil {
		pool.Close()
		log.Error().Err(err).Msg("failed to auto-migrate database")
		return nil, fmt.Errorf("auto-migrate database: %w", err)
	}

	return &Database{
		pool: pool,
		log:  log,
	}, nil
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) HealthCheckWithContext(ctx context.Context) error {
	if err := db.pool.Ping(ctx); err != nil {
		db.log.Error().Err(err).Msg("database health check failed")
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *Database) Shutdown(ctx context.Context) error {
	if db.pool != nil {
		db.log.Info().Msg("closing postgres pool")
		db.pool.Close()
	}

	return nil
}
