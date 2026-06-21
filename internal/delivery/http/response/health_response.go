package response

import (
	"time"

	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
)

type HealthResponse struct {
	Name        string    `json:"name" example:"BE Plant Factory"`
	Environment string    `json:"environment" example:"production"`
	Status      string    `json:"status" example:"ok"`
	Timestamp   time.Time `json:"timestamp" example:"2026-01-01T00:00:00Z"`
}

func NewHealthResponse(dto health.StatusResponse) HealthResponse {
	return HealthResponse{
		Name:        dto.Name,
		Environment: dto.Environment,
		Status:      dto.Status,
		Timestamp:   dto.Timestamp,
	}
}
