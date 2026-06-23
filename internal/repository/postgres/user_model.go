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
	LastLoggedInAt time.Time          `db:"last_logged_in_at"`
	CreatedAt      time.Time          `db:"created_at"`
	UpdatedAt      *time.Time         `db:"updated_at"`
}

type UpsertSocialUser struct {
	ID          uuid.UUID          `db:"id"`
	CurrentStep entity.CurrentStep `db:"current_step"`
}
