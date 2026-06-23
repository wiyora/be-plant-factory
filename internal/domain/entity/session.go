package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	RefreshTokenHash string
	DeviceName       string
	IPAddress        string
	ExpiredAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}
