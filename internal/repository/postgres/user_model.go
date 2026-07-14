package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

type User struct {
	ID             uuid.UUID          `db:"id"`
	Email          string             `db:"email"`
	Name           string             `db:"name"`
	Avatar         string             `db:"avatar"`
	IsSuperAdmin   bool               `db:"is_super_admin"`
	CurrentStep    entity.CurrentStep `db:"current_step"`
	Status         entity.UserStatus  `db:"status"`
	LastLoggedInAt time.Time          `db:"last_logged_in_at"`
	CreatedAt      time.Time          `db:"created_at"`
	UpdatedAt      *time.Time         `db:"updated_at"`
}

func (u User) ToEntity() entity.User {
	return entity.User{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name,
		Avatar:         u.Avatar,
		IsSuperAdmin:   u.IsSuperAdmin,
		CurrentStep:    u.CurrentStep,
		Status:         u.Status,
		LastLoggedInAt: u.LastLoggedInAt,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

type ListUser struct {
	ID        uuid.UUID         `db:"id"`
	Email     string            `db:"email"`
	Name      string            `db:"name"`
	Avatar    string            `db:"avatar"`
	Status    entity.UserStatus `db:"status"`
	CreatedAt time.Time         `db:"created_at"`
}

func (u ListUser) ToEntity() entity.User {
	return entity.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

type UpsertSocialUser struct {
	ID          uuid.UUID          `db:"id"`
	CurrentStep entity.CurrentStep `db:"current_step"`
	Status      entity.UserStatus  `db:"status"`
}

type DropdownUser struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (u DropdownUser) ToEntity() entity.DropdownItem {
	return entity.DropdownItem{
		ID:   u.ID,
		Name: u.Name,
	}
}
