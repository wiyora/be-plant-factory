package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserTenant struct {
	UserID    uuid.UUID
	TenantID  uuid.UUID
	RoleID    uuid.UUID
	User      User
	Tenant    Tenant
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type TenantUserFilter struct {
	TenantID   uuid.UUID
	Search     Search
	Pagination Pagination
	Order      Order
}

type UserTenantFilter struct {
	UserID     uuid.UUID
	TenantIDs  []uuid.UUID
	RoleIDs    []uuid.UUID
	Pagination Pagination
	Order      Order
}

type TenantUserList struct {
	UserID       uuid.UUID
	UserName     string
	UserEmail    string
	UserAvatar   string
	UserStatus   UserStatus
	RoleID       uuid.UUID
	RoleName     string
	AssignedDate time.Time
}

type UserTenantList struct {
	UserTenantID uuid.UUID
	TenantID     uuid.UUID
	TenantName   string
	TenantLogo   string
	RoleID       uuid.UUID
	RoleName     string
	AssignedDate time.Time
}
