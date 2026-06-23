package auth

import (
	"time"

	"github.com/google/uuid"
)

type GoogleRequest struct {
	RedirectUri string
	BaseUrl     string
}

type GenerateRefreshToken struct {
	RefreshToken     string
	HashRefreshToken string
	ExpiresAt        time.Time
}

type GenerateAccessToken struct {
	AccessToken string
	ExpiresAt   time.Time
}

type CallbackRequest struct {
	State      string
	Code       string
	DeviceName string
	IPAddress  string
}

type CallbackResponse struct {
	TokenResponse
	BaseURL     string
	RedirectUri string
	IsCompleted bool
}

type TokenResponse struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type LogoutRequest struct {
	AccessTokenID uuid.UUID
	UserID        uuid.UUID
	SessionID     uuid.UUID
	ExpiredAt     time.Time
}
