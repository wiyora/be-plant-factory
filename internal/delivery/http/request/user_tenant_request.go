package request

import (
	"github.com/google/uuid"
)

type AssignUserRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}
