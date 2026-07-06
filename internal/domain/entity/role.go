package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID              uuid.UUID
	Name            string
	TotalPermission int
	TotalUser       int
	Permissions     []Permission
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type RoleFilter struct {
	Search     Search
	Pagination Pagination
	Order      Order
}
