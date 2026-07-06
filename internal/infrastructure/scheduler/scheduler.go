package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

type Scheduler interface {
	Scheduler() gocron.Scheduler
	Log() zerolog.Logger
	Start()
	Shutdown(ctx context.Context) error
}

type scheduler struct {
	log       zerolog.Logger
	scheduler gocron.Scheduler
}

func New(i do.Injector) (Scheduler, error) {
	rawLog := do.MustInvoke[*zerolog.Logger](i)

	log := logger.WithLayer(rawLog, logger.LayerScheduler)
	log.Info().Msg("Initializing scheduler")

	s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		log.Error().Err(err).Msg("failed to create scheduler")
		return nil, err
	}

	return scheduler{
		log:       log,
		scheduler: s,
	}, nil
}

func (s scheduler) Start() {
	s.log.Info().Msg("Starting scheduler")
	s.scheduler.Start()
}

func (s scheduler) Scheduler() gocron.Scheduler {
	return s.scheduler
}

func (s scheduler) Log() zerolog.Logger {
	return s.log
}

func (s scheduler) Shutdown(ctx context.Context) error {
	s.log.Info().Msg("Shutting down scheduler")

	if err := s.scheduler.Shutdown(); err != nil {
		s.log.Error().Err(err).Msg("failed to shutdown scheduler")
		return err
	}

	return nil
}
