package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/jwt"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/cache"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/redis"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Middleware interface {
	Register(app *fiber.App)
	Auth() fiber.Handler
}

type middleware struct {
	conf  *config.Config        `do:""`
	log   *zerolog.Logger       `do:""`
	jwt   jwt.Service           `do:""`
	cache cache.UserRepository  `do:""`
	token redis.TokenRepository `do:""`
}

func New(i do.Injector) (Middleware, error) {
	return do.InvokeStruct[*middleware](i)
}
