package bootstrap

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/cron"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/handler"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/middleware"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/route"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/websocket"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure"

	"github.com/rizalarfiyan/be-plant-factory/internal/repository"
	useCase "github.com/rizalarfiyan/be-plant-factory/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

func NewContainer(log *zerolog.Logger, conf *config.Config) (do.Injector, error) {
	injector := do.New(
		infrastructure.Package,
		repository.Package,
		useCase.Package,
		handler.Package,
	)

	// General
	do.ProvideValue(injector, context.Background())
	do.ProvideValue(injector, log)
	do.ProvideValue(injector, conf)

	// HTTP
	do.Provide(injector, route.New)
	do.Provide(injector, middleware.New)
	do.Provide(injector, NewServer)

	// Cron
	do.Provide(injector, cron.New)

	// WebSocket
	do.Provide(injector, websocket.NewSocketHandler)

	return injector, nil
}
