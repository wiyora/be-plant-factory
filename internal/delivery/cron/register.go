package cron

import "context"

type CronFunc func() error

type CronTask func(ctx context.Context) func()

func (c cron) jobs() []CronFunc {
	return []CronFunc{
		c.registerDailyJob("storage-cleanup-temporary", 0, 15, c.storage.CleanupTemporary),
	}
}

func (c cron) MapJobs() error {
	for _, job := range c.jobs() {
		if err := job(); err != nil {
			return err
		}
	}

	return nil
}
