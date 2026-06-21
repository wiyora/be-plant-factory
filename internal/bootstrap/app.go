package bootstrap

import (
	"fmt"

	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	loggerInfra "github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type App struct {
	server *Server
	inject do.Injector
	log    zerolog.Logger
}

func NewApp() (*App, error) {
	env := config.MustLoadAppEnv()
	logger, err := loggerInfra.New(env)
	if err != nil {
		return nil, err
	}

	logConf := loggerInfra.WithLayer(logger, loggerInfra.LayerConfig)
	conf, err := config.NewConfig(".env", logConf).Load()
	if err != nil {
		return nil, err
	}

	level, err := zerolog.ParseLevel(conf.Log.Level)
	if err != nil {
		logConf.Error().Err(err).Str("level", conf.Log.Level).Msg("invalid log level, defaulting to info")
		return nil, err
	}

	zerolog.SetGlobalLevel(level)

	injector, err := NewContainer(logger, conf)
	if err != nil {
		return nil, err
	}

	server, err := do.Invoke[*Server](injector)
	if err != nil {
		return nil, err
	}

	appLog := loggerInfra.WithLayer(logger, loggerInfra.LayerApp)
	return &App{
		server: server,
		inject: injector,
		log:    appLog,
	}, nil
}

func (a *App) Run() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			_ = a.inject.RootScope().Shutdown()
		}
	}()

	err = a.server.Start()
	if err != nil {
		a.log.Error().Err(err).Msg("failed to start http server")
		return err
	}

	_, report := a.inject.RootScope().ShutdownOnSignals()
	if report.Succeed {
		a.log.Info().Msgf("app shutdown gracefully in %v", report.ShutdownTime)

		return nil
	}

	a.log.Error().Str("detail", report.Error()).Msg("app shutdown with error")
	return fmt.Errorf("app shutdown failed")
}
