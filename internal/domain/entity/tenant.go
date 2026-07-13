package entity

import (
	"time"

	"github.com/google/uuid"
)

// Tenant Status
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
)

func (s TenantStatus) String() string {
	return string(s)
}

func (s TenantStatus) Valid() bool {
	switch s {
	case TenantStatusActive, TenantStatusInactive, TenantStatusSuspended:
		return true
	default:
		return false
	}
}

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Logo      string
	Status    TenantStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type TenantFilter struct {
	Search     Search
	Pagination Pagination
	Order      Order
	Status     TenantStatus
}
