package request

import (
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type CreateUpdateRoleRequest struct {
	ID          uuid.UUID           `json:"id" saggerignore:"true"`
	Name        string              `json:"name" validate:"required,min=3,max=32" example:"Admin"`
	Permissions []entity.Permission `json:"permissions" validate:"required,min=1,unique,dive,enum" example:"user:list,user:detail"`
}

func (r CreateUpdateRoleRequest) ToEntity() entity.Role {
	return entity.Role{
		ID:              r.ID,
		Name:            r.Name,
		Permissions:     r.Permissions,
		TotalPermission: len(r.Permissions),
	}
}
