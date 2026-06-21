package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

type AccessClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

func (s service) GenerateAccessToken(userID, sessionID uuid.UUID, issuedAt time.Time) (string, error) {
	jwtID, err := helper.ID()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT ID: %w", err)
	}

	claims := AccessClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jwtID.String(),
			Subject:   userID.String(),
			Issuer:    s.conf.JWT.Issuer,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(s.conf.Auth.AccessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.conf.JWT.Secret))
}

func (s service) ParseAccessToken(tokenString string) (AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}

		return []byte(s.conf.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return AccessClaims{}, fmt.Errorf("parse access token: %w", err)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid || claims == nil {
		return AccessClaims{}, errors.New("invalid access token")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) {
		return AccessClaims{}, jwt.ErrTokenExpired
	}

	return *claims, nil
}
