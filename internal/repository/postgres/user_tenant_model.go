package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type ListTenantUser struct {
	UserID       uuid.UUID         `db:"user_id"`
	UserName     string            `db:"user_name"`
	UserEmail    string            `db:"user_email"`
	UserAvatar   string            `db:"user_avatar"`
	UserStatus   entity.UserStatus `db:"user_status"`
	RoleID       uuid.UUID         `db:"role_id"`
	RoleName     string            `db:"role_name"`
	AssignedDate time.Time         `db:"assigned_date"`
}

func (m ListTenantUser) ToEntity() entity.TenantUserList {
	return entity.TenantUserList{
		UserID:       m.UserID,
		UserName:     m.UserName,
		UserEmail:    m.UserEmail,
		UserAvatar:   m.UserAvatar,
		UserStatus:   m.UserStatus,
		RoleID:       m.RoleID,
		RoleName:     m.RoleName,
		AssignedDate: m.AssignedDate,
	}
}

type ListUserTenant struct {
	UserTenantID uuid.UUID `db:"user_tenant_id"`
	TenantID     uuid.UUID `db:"tenant_id"`
	TenantName   string    `db:"tenant_name"`
	TenantLogo   string    `db:"tenant_logo"`
	RoleID       uuid.UUID `db:"role_id"`
	RoleName     string    `db:"role_name"`
	AssignedDate time.Time `db:"assigned_date"`
}

func (m ListUserTenant) ToEntity() entity.UserTenantList {
	return entity.UserTenantList{
		UserTenantID: m.UserTenantID,
		TenantID:     m.TenantID,
		TenantName:   m.TenantName,
		TenantLogo:   m.TenantLogo,
		RoleID:       m.RoleID,
		RoleName:     m.RoleName,
		AssignedDate: m.AssignedDate,
	}
}
