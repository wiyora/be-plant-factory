package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/handler"
	"github.com/samber/do/v2"
)

type Router interface {
	Register(app *fiber.App)
}

type router struct {
	health handler.HealthHandler `do:""`
}

func New(i do.Injector) (Router, error) {
	return do.InvokeStruct[*router](i)
}

func (r *router) Register(app *fiber.App) {
	app.Get("/health", r.health.Check)
}
