package entity

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Permission string

const (
	PermissionSuperAdmin Permission = "super-admin"

	PermissionUserList             Permission = "user:list"
	PermissionUserDetail           Permission = "user:detail"
	PermissionUserCreate           Permission = "user:create"
	PermissionUserUpdate           Permission = "user:update"
	PermissionUserUpdateStatus     Permission = "user:update-status"
	PermissionUserAssignRoleTenant Permission = "user:assign-role-tenant"

	PermissionRoleList   Permission = "role:list"
	PermissionRoleDetail Permission = "role:detail"
	PermissionRoleCreate Permission = "role:create"
	PermissionRoleUpdate Permission = "role:update"
	PermissionRoleDelete Permission = "role:delete"

	PermissionTenantList         Permission = "tenant:list"
	PermissionTenantDetail       Permission = "tenant:detail"
	PermissionTenantCreate       Permission = "tenant:create"
	PermissionTenantUpdate       Permission = "tenant:update"
	PermissionTenantUpdateStatus Permission = "tenant:update-status"
)

var MapPermissions = mapset.NewSet(
	PermissionSuperAdmin,

	PermissionUserList,
	PermissionUserDetail,
	PermissionUserCreate,
	PermissionUserUpdate,
	PermissionUserUpdateStatus,
	PermissionUserAssignRoleTenant,

	PermissionRoleList,
	PermissionRoleDetail,
	PermissionRoleCreate,
	PermissionRoleUpdate,
	PermissionRoleDelete,

	PermissionTenantList,
	PermissionTenantDetail,
	PermissionTenantCreate,
	PermissionTenantUpdate,
	PermissionTenantUpdateStatus,
)

func (p Permission) Valid() bool {
	return MapPermissions.Contains(p)
}

func (p Permission) IsSuperAdmin() bool {
	return p == PermissionSuperAdmin
}

func (p Permission) String() string {
	return string(p)
}
