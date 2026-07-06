package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type ListRoleResponse struct {
	ID              uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name            string    `json:"name" example:"Admin"`
	TotalPermission int       `json:"total_permission" example:"5"`
	TotalUser       int       `json:"total_user" example:"10"`
}

func NewListRoleResponse(r entity.Role) ListRoleResponse {
	return ListRoleResponse{
		ID:              r.ID,
		Name:            r.Name,
		TotalPermission: r.TotalPermission,
		TotalUser:       r.TotalUser,
	}
}

type DetailRoleResponse struct {
	ID              uuid.UUID           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name            string              `json:"name" example:"Admin"`
	TotalPermission int                 `json:"total_permission" example:"5"`
	Permissions     []entity.Permission `json:"permissions"`
	CreatedAt       time.Time           `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       *time.Time          `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewDetailRoleResponse(r entity.Role) DetailRoleResponse {
	return DetailRoleResponse{
		ID:              r.ID,
		Name:            r.Name,
		TotalPermission: r.TotalPermission,
		Permissions:     r.Permissions,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
