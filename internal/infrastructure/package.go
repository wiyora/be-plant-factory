package infrastructure

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/jwt"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/redis"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/validator"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(postgres.New),
	do.Lazy(redis.New),
	do.Lazy(validator.New),
	do.Lazy(jwt.New),
)
