package cron

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (c cron) registerDailyJob(name string, hour, minute uint, task CronTask) CronFunc {
	log := c.scheduler.Log().With().Str("name", name).Uint("hour", hour).Uint("minute", minute).Logger()
	ctx := logger.WithSection(&log, logger.SectionCron).WithContext(context.Background())

	return func() error {
		_, err := c.scheduler.Scheduler().NewJob(
			gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(hour, minute, 0))),
			gocron.NewTask(task(ctx)),
			gocron.WithName(name),
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to register daily job")
			return err
		}

		log.Info().Msg("daily job registered")
		return nil
	}
}
