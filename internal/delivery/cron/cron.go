package cron

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/scheduler"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/storage"
	"github.com/samber/do/v2"
)

type Cron interface {
	MapJobs() error
}

type cron struct {
	scheduler scheduler.Scheduler    `do:""`
	storage   storage.StorageUseCase `do:""`
}

func New(i do.Injector) (Cron, error) {
	return do.InvokeStruct[cron](i)
}
