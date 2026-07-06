package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type ListTenantResponse struct {
	ID     uuid.UUID           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name   string              `json:"name" example:"My Tenant"`
	Logo   string              `json:"logo" example:"tenant/logo/uuid.png"`
	Status entity.TenantStatus `json:"status" example:"active"`
}

func NewListTenantResponse(t entity.Tenant) ListTenantResponse {
	return ListTenantResponse{
		ID:     t.ID,
		Name:   t.Name,
		Logo:   t.Logo,
		Status: t.Status,
	}
}

type DetailTenantResponse struct {
	ID        uuid.UUID           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string              `json:"name" example:"My Tenant"`
	Logo      string              `json:"logo" example:"tenant/logo/uuid.png"`
	Status    entity.TenantStatus `json:"status" example:"active"`
	CreatedAt time.Time           `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty" example:"2024-01-01T00:00:00Z"`
}

func NewDetailTenantResponse(t entity.Tenant) DetailTenantResponse {
	return DetailTenantResponse{
		ID:        t.ID,
		Name:      t.Name,
		Logo:      t.Logo,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
