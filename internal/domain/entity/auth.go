package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthContext struct {
	ID          uuid.UUID
	Email       string
	Name        string
	Avatar      string
	CurrentStep CurrentStep
}

type AccessTokenContext struct {
	ID        uuid.UUID
	ExpiredAt time.Time
}
