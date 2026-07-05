package entity

import (
	"time"

	"github.com/google/uuid"
)

// User Current Step
type CurrentStep string

const (
	CurrentStepInitial   CurrentStep = "initial"
	CurrentStepCompleted CurrentStep = "completed"
)

func (s CurrentStep) String() string {
	return string(s)
}

func (s CurrentStep) Valid() bool {
	switch s {
	case CurrentStepInitial, CurrentStepCompleted:
		return true
	default:
		return false
	}
}

// User Status
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) Valid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusBanned:
		return true
	default:
		return false
	}
}

type User struct {
	ID             uuid.UUID
	Email          string
	Name           string
	Avatar         string
	IsSuperAdmin   bool
	CurrentStep    CurrentStep
	Status         UserStatus
	LastLoggedInAt time.Time
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type UserFilter struct {
	Search     Search
	Pagination Pagination
	Order      Order
	Status     UserStatus
}
