package health

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/samber/do/v2"
)

type HealthUseCase interface {
	Check() StatusResponse
}

type healthUseCase struct {
	conf *config.Config `do:""`
}

func NewHealthUseCase(i do.Injector) (HealthUseCase, error) {
	return do.InvokeStruct[healthUseCase](i)
}
