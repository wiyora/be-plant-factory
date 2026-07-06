package response

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

type RolePermission struct {
	Key     entity.Permission   `json:"key" example:"user:create"`
	Label   string              `json:"label" example:"Create User"`
	Require []entity.Permission `json:"require,omitempty" example:"user:list"`
}

type PermissionModule struct {
	Name        string           `json:"name" example:"User Management"`
	Permissions []RolePermission `json:"permissions"`
}

func NewPermissionModule(r entity.PermissionModule) PermissionModule {
	permissions := make([]RolePermission, len(r.Permissions))
	for i, p := range r.Permissions {
		permissions[i] = RolePermission{
			Key:     p.Key,
			Label:   p.Label,
			Require: p.Require,
		}
	}

	return PermissionModule{
		Name:        r.Name,
		Permissions: permissions,
	}
}
