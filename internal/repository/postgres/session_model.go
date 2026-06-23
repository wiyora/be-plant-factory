package postgres

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID  `db:"id"`
	UserID           uuid.UUID  `db:"user_id"`
	RefreshTokenHash string     `db:"refresh_token_hash"`
	DeviceName       string     `db:"device_name"`
	IPAddress        netip.Addr `db:"ip_address"`
	ExpiredAt        time.Time  `db:"expired_at"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}
