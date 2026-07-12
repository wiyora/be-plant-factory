package bootstrap

import (
	"context"
	"errors"
	"net"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/cron"
	httpError "github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/middleware"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/http/route"
	"github.com/rizalarfiyan/be-plant-factory/internal/delivery/websocket"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/scheduler"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/validator"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Server struct {
	app       *fiber.App
	conf      *config.Config
	log       zerolog.Logger
	route     route.Router
	scheduler scheduler.Scheduler
}

func NewServer(i do.Injector) (*Server, error) {
	conf := do.MustInvoke[*config.Config](i)
	rawLog := do.MustInvoke[*zerolog.Logger](i)
	mid := do.MustInvoke[middleware.Middleware](i)
	route := do.MustInvoke[route.Router](i)
	validate := do.MustInvoke[*validator.Validate](i)
	scheduler := do.MustInvoke[scheduler.Scheduler](i)
	cron := do.MustInvoke[cron.Cron](i)
	socketHandler := do.MustInvoke[websocket.SocketHandler](i)

	log := logger.WithLayer(rawLog, logger.LayerHttp)

	if err := cron.MapJobs(); err != nil {
		log.Error().Err(err).Msg("Failed to map cron jobs")
		return nil, err
	}

	log.Info().Msg("Initializing http server")

	validate.Register()

	if err := validate.CustomValidation(); err != nil {
		log.Error().Err(err).Msg("Failed to register custom validation")
		return nil, err
	}

	app := fiber.New(fiber.Config{
		AppName:             conf.App.Name,
		BodyLimit:           conf.Fiber.BodyLimit,
		ErrorHandler:        httpError.Handle,
		JSONEncoder:         sonic.Marshal,
		JSONDecoder:         sonic.Unmarshal,
		StructValidator:     validate,
		PassLocalsToContext: true,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: conf.Fiber.TrustedProxies,
		},
	})

	app.State().Set(constant.AppStateValidatorKey, validate.Validator())

	log.Info().Msg("Registering http server middleware")
	mid.Register(app)

	log.Info().Msg("Registering http server routes")
	route.Register(app)

	log.Info().Msg("Registering websocket routes")
	websocket.Register(app, socketHandler)

	return &Server{
		conf:      conf,
		log:       log,
		route:     route,
		app:       app,
		scheduler: scheduler,
	}, nil
}

func (s *Server) Start() (err error) {
	s.scheduler.Start()

	s.app.Server().ReadTimeout = s.conf.HTTP.ReadTimeout
	s.app.Server().WriteTimeout = s.conf.HTTP.WriteTimeout
	s.app.Server().IdleTimeout = s.conf.HTTP.IdleTimeout

	go func() {
		s.log.Info().Str("address", s.conf.HTTP.AddressWithScheme()).Msg("starting http server")

		config := fiber.ListenConfig{
			EnablePrefork:         s.conf.Fiber.Prefork,
			DisableStartupMessage: s.conf.App.Env.IsServerEnv(),
			ShutdownTimeout:       s.conf.HTTP.ShutdownTimeout,
		}

		if err := s.app.Listen(s.conf.HTTP.Address(), config); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info().Msg("Stopping http server")
	return s.app.ShutdownWithContext(ctx)
}
