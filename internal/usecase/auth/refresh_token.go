package auth

import (
	"context"
	"time"

	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u authUseCase) RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	hashRefreshToken := helper.SHA256Hex(refreshToken)
	currentRt, err := u.sessionRepo.GetByRefreshToken(ctx, hashRefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to get session by refresh token")
		return TokenResponse{}, err
	}

	if helper.IsEmptyStruct(currentRt) {
		log.Warn().Msg("refresh token not found")
		return TokenResponse{}, domainError.New(code.InvalidRefreshToken)
	}

	generateTime := currentRt.ExpiredAt.Add(-u.conf.Auth.RefreshTokenDuration)
	minimalRetry := generateTime.Add(u.conf.Auth.AccessTokenDuration - u.conf.Auth.ToleranceDuration)
	remainingTime := time.Since(minimalRetry)

	if remainingTime < 0 {
		log.Warn().Msg("refresh token is not eligible for refresh")
		return TokenResponse{}, domainError.New(code.RefreshTokenNotEligible, domainError.WithParams(map[string]any{
			"remaining_seconds": helper.GetRemainingSeconds(remainingTime),
		}))
	}

	rt, err := u.generateRefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
		return TokenResponse{}, err
	}

	at, err := u.generateAccessToken(currentRt.UserID, currentRt.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		return TokenResponse{}, err
	}

	if err := u.sessionRepo.Update(ctx, currentRt.ID, rt.HashRefreshToken, rt.ExpiresAt); err != nil {
		log.Error().Err(err).Msg("failed to update refresh token")
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:           at.AccessToken,
		AccessTokenExpiresAt:  at.ExpiresAt,
		RefreshToken:          rt.RefreshToken,
		RefreshTokenExpiresAt: rt.ExpiresAt,
	}, nil
}
