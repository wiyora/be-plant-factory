package health

import (
	"time"
)

func (uc healthUseCase) Check() StatusResponse {
	return StatusResponse{
		Name:        uc.conf.App.Name,
		Environment: uc.conf.App.Env.String(),
		Status:      "ok",
		Timestamp:   time.Now().UTC(),
	}
}
