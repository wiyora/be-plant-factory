package auth

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u authUseCase) generateRefreshToken() (GenerateRefreshToken, error) {
	refreshToken, err := helper.GenerateRandomToken(32)
	if err != nil {
		return GenerateRefreshToken{}, err
	}

	return GenerateRefreshToken{
		RefreshToken:     refreshToken,
		HashRefreshToken: helper.SHA256Hex(refreshToken),
		ExpiresAt:        time.Now().Add(u.conf.Auth.RefreshTokenDuration),
	}, nil
}

func (u authUseCase) generateAccessToken(userId, sessionId uuid.UUID) (GenerateAccessToken, error) {
	accessToken, err := u.jwtService.GenerateAccessToken(userId, sessionId, time.Now())
	if err != nil {
		return GenerateAccessToken{}, err
	}

	return GenerateAccessToken{
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(u.conf.Auth.AccessTokenDuration),
	}, nil
}

func (u authUseCase) GetDeviceName(rawUA string) string {
	p := useragent.Parse(rawUA)
	ua := p.Name
	if p.Version != "" {
		ua += " v" + p.Version
	}

	if p.OS != "" {
		ua += " on " + p.OS
	}

	if ua == "" {
		return "Unknown"
	}

	if p.Mobile {
		ua = "(Mobile) " + ua
	}
	if p.Tablet {
		ua = "(Tablet) " + ua
	}
	if p.Desktop {
		ua = "(Desktop) " + ua
	}
	if p.Bot {
		ua = "(Bot) " + ua
	}

	return strings.TrimSpace(ua)
}
