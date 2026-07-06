package repository

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/cache"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/redis"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/storage"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	// Generic Manager
	do.Lazy(postgres.NewTransactionManager),

	// Postgres repositories
	do.Lazy(postgres.NewUserRepository),
	do.Lazy(postgres.NewSessionRepository),
	do.Lazy(postgres.NewRoleRepository),

	// Redis repositories
	do.Lazy(redis.NewCacheRepository),
	do.Lazy(redis.NewTokenRepository),
	do.Lazy(redis.NewStateRepository),

	// Cache repositories
	do.Lazy(cache.NewUserRepository),

	// Storage repositories
	do.Lazy(storage.NewS3Repository),
)
