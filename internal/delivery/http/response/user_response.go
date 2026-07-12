package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type ListUserResponse struct {
	ID        uuid.UUID         `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string            `json:"email" example:"rizalarfiyan@plant-factory.com"`
	Name      string            `json:"name" example:"Rizal Arfiyan"`
	Avatar    string            `json:"avatar" example:"avatar: coding"`
	Status    entity.UserStatus `json:"status" example:"active"`
	CreatedAt time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

func NewListUserResponse(u entity.User) ListUserResponse {
	return ListUserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

type DetailUserResponse struct {
	ID             uuid.UUID          `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email          string             `json:"email" example:"rizalarfiyan@plant-factory.com"`
	Name           string             `json:"name" example:"Rizal Arfiyan"`
	Avatar         string             `json:"avatar" example:"avatar: coding"`
	CurrentStep    entity.CurrentStep `json:"current_step" example:"completed"`
	Status         entity.UserStatus  `json:"status" example:"active"`
	IsSuperAdmin   bool               `json:"is_super_admin" example:"false"`
	LastLoggedInAt time.Time          `json:"last_logged_in_at" example:"2024-01-01T00:00:00Z"`
	CreatedAt      time.Time          `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty" example:"2024-01-01T00:00:00Z"`
}

func NewDetailUserResponse(u entity.User) DetailUserResponse {
	return DetailUserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name,
		Avatar:         u.Avatar,
		CurrentStep:    u.CurrentStep,
		Status:         u.Status,
		IsSuperAdmin:   u.IsSuperAdmin,
		LastLoggedInAt: u.LastLoggedInAt,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
