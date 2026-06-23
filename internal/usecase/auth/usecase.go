package auth

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/jwt"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/social"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/redis"
	"github.com/samber/do/v2"
)

type AuthUseCase interface {
	GoogleLoginURL(ctx context.Context, req GoogleRequest) (string, error)
	HandleGoogleCallback(ctx context.Context, req CallbackRequest) (CallbackResponse, error)
	GetDeviceName(rawUA string) string
	RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error)
	Logout(ctx context.Context, req LogoutRequest) error
}

type authUseCase struct {
	conf        *config.Config             `do:""`
	authManager social.Manager             `do:""`
	jwtService  jwt.Service                `do:""`
	stateRepo   redis.StateRepository      `do:""`
	tokenRepo   redis.TokenRepository      `do:""`
	userRepo    postgres.UserRepository    `do:""`
	sessionRepo postgres.SessionRepository `do:""`
}

func NewAuthUseCase(i do.Injector) (AuthUseCase, error) {
	return do.InvokeStruct[authUseCase](i)
}
