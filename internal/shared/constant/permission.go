package constant

import "github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"

var PermissionModules = []entity.PermissionModule{
	{
		Name: "User",
		Permissions: []entity.RolePermission{
			{
				Key:   entity.PermissionUserList,
				Label: "User List",
			},
			{
				Key:   entity.PermissionUserDetail,
				Label: "User Detail",
			},
			{
				Key:   entity.PermissionUserCreate,
				Label: "User Create",
				Require: []entity.Permission{
					entity.PermissionUserList,
				},
			},
			{
				Key:   entity.PermissionUserUpdate,
				Label: "User Update",
				Require: []entity.Permission{
					entity.PermissionUserList,
				},
			},
			{
				Key:   entity.PermissionUserUpdateStatus,
				Label: "User Update Status",
				Require: []entity.Permission{
					entity.PermissionUserList,
				},
			},
			{
				Key:   entity.PermissionUserAssignRoleTenant,
				Label: "User Assign Role Tenant",
			},
		},
	},
	{
		Name: "Role",
		Permissions: []entity.RolePermission{
			{
				Key:   entity.PermissionRoleList,
				Label: "Role List",
			},
			{
				Key:   entity.PermissionRoleDetail,
				Label: "Role Detail",
			},
			{
				Key:   entity.PermissionRoleCreate,
				Label: "Role Create",
				Require: []entity.Permission{
					entity.PermissionRoleList,
				},
			},
			{
				Key:   entity.PermissionRoleUpdate,
				Label: "Role Update",
				Require: []entity.Permission{
					entity.PermissionRoleList,
				},
			},
			{
				Key:   entity.PermissionRoleDelete,
				Label: "Role Delete",
				Require: []entity.Permission{
					entity.PermissionRoleList,
				},
			},
		},
	},
	{
		Name: "Tenant",
		Permissions: []entity.RolePermission{
			{
				Key:   entity.PermissionTenantList,
				Label: "Tenant List",
			},
			{
				Key:   entity.PermissionTenantDetail,
				Label: "Tenant Detail",
			},
			{
				Key:   entity.PermissionTenantCreate,
				Label: "Tenant Create",
				Require: []entity.Permission{
					entity.PermissionTenantList,
				},
			},
			{
				Key:   entity.PermissionTenantUpdate,
				Label: "Tenant Update",
				Require: []entity.Permission{
					entity.PermissionTenantList,
				},
			},
			{
				Key:   entity.PermissionTenantUpdateStatus,
				Label: "Tenant Update Status",
				Require: []entity.Permission{
					entity.PermissionTenantList,
				},
			},
		},
	},
}
