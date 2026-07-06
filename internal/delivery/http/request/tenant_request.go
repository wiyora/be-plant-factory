package request

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type CreateTenantRequest struct {
	Name string `json:"name" validate:"required,alphaspace,min=3,max=32" example:"My Tenant"`
	Logo string `json:"logo" validate:"required" example:"temp/tenant/logo/123e4567-e89b-12d3-a456-426614174000.png"`
}

func (r CreateTenantRequest) ToEntity() entity.Tenant {
	return entity.Tenant{
		Name: r.Name,
		Logo: r.Logo,
	}
}

type UpdateTenantRequest struct {
	Name string `json:"name" validate:"required,alphaspace,min=3,max=32" example:"My Tenant"`
	Logo string `json:"logo" validate:"omitempty" example:"temp/tenant/logo/123e4567-e89b-12d3-a456-426614174000.png"`
}

func (r UpdateTenantRequest) ToEntity() entity.Tenant {
	return entity.Tenant{
		Name: r.Name,
		Logo: r.Logo,
	}
}
