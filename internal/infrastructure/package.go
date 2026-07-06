package infrastructure

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/jwt"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/redis"
	s3client "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/s3client"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/scheduler"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/social"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/validator"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(postgres.New),
	do.Lazy(redis.New),
	do.Lazy(validator.New),
	do.Lazy(jwt.New),
	do.Lazy(social.New),
	do.Lazy(s3client.New),
	do.Lazy(scheduler.New),
)
