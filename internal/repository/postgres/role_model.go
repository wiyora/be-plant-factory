package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type Role struct {
	ID              uuid.UUID           `db:"id"`
	Name            string              `db:"name"`
	TotalPermission int                 `db:"total_permission"`
	Permissions     []entity.Permission `db:"permissions"`
	CreatedAt       time.Time           `db:"created_at"`
	UpdatedAt       *time.Time          `db:"updated_at"`
}

func (r Role) ToEntity() entity.Role {
	return entity.Role{
		ID:              r.ID,
		Name:            r.Name,
		TotalPermission: r.TotalPermission,
		Permissions:     r.Permissions,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func (r ListRole) ToEntity() entity.Role {
	return entity.Role{
		ID:              r.ID,
		Name:            r.Name,
		TotalPermission: r.TotalPermission,
		TotalUser:       r.TotalUser,
	}
}

type ListRole struct {
	ID              uuid.UUID `db:"id"`
	Name            string    `db:"name"`
	TotalPermission int       `db:"total_permission"`
	TotalUser       int       `db:"total_user"`
}
