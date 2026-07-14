package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type ListTenantUserResponse struct {
	ID           uuid.UUID         `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Avatar       string            `json:"avatar" example:"avatar:coding"`
	Email        string            `json:"email" example:"user@example.com"`
	Name         string            `json:"name" example:"John Doe"`
	Status       entity.UserStatus `json:"status" example:"active"`
	Role         RoleResponse      `json:"role"`
	AssignedDate time.Time         `json:"assigned_date" example:"2024-01-01T00:00:00Z"`
}

type RoleResponse struct {
	ID   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string    `json:"name" example:"Admin"`
}

func NewListTenantUserResponse(u entity.TenantUserList) ListTenantUserResponse {
	return ListTenantUserResponse{
		ID:     u.UserID,
		Avatar: u.UserAvatar,
		Email:  u.UserEmail,
		Name:   u.UserName,
		Status: u.UserStatus,
		Role: RoleResponse{
			ID:   u.RoleID,
			Name: u.RoleName,
		},
		AssignedDate: u.AssignedDate,
	}
}

type ListUserTenantResponse struct {
	ID           uuid.UUID              `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tenant       TenantUserListResponse `json:"tenant"`
	Role         RoleResponse           `json:"role"`
	AssignedDate time.Time              `json:"assigned_date" example:"2024-01-01T00:00:00Z"`
}

type TenantUserListResponse struct {
	ID   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string    `json:"name" example:"Tenant A"`
	Logo string    `json:"logo" example:"tenant:logo123"`
}

func NewListUserTenantResponse(t entity.UserTenantList) ListUserTenantResponse {
	return ListUserTenantResponse{
		ID: t.UserTenantID,
		Tenant: TenantUserListResponse{
			ID:   t.TenantID,
			Name: t.TenantName,
			Logo: t.TenantLogo,
		},
		Role: RoleResponse{
			ID:   t.RoleID,
			Name: t.RoleName,
		},
		AssignedDate: t.AssignedDate,
	}
}
