package entity

import (
	"time"

	"github.com/google/uuid"
)

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

type User struct {
	ID             uuid.UUID
	Email          string
	Name           string
	Avatar         string
	IsSuperAdmin   bool
	CurrentStep    CurrentStep
	LastLoggedInAt time.Time
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
