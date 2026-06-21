package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Middleware interface {
	Register(app *fiber.App)
}

type middleware struct {
	conf *config.Config  `do:""`
	log  *zerolog.Logger `do:""`
}

func New(i do.Injector) (Middleware, error) {
	return do.InvokeStruct[*middleware](i)
}
