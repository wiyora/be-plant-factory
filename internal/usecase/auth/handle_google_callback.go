package auth

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u authUseCase) HandleGoogleCallback(ctx context.Context, req CallbackRequest) (CallbackResponse, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)
	providerName := entity.ProviderGoogle

	res := CallbackResponse{
		BaseURL: u.conf.Frontend.DefaultAllowedUrl(),
	}

	value, err := u.stateRepo.ValidateAndDelete(ctx, providerName, req.State)
	if err != nil || helper.IsEmptyStruct(value) {
		log.Error().Err(err).Msg("failed to validate state")
		return res, domainError.New(code.InvalidState)
	}

	res.BaseURL = value.BaseURL
	res.RedirectUri = value.RedirectUri

	provider, err := u.authManager.GetProvider(entity.ProviderGoogle)
	if err != nil {
		log.Error().Err(err).Msg("failed to get auth provider")
		return res, domainError.New(code.ProviderError)
	}

	socialUser, err := provider.GetUserInfo(ctx, req.Code)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info from provider")
		return res, domainError.New(code.TokenExchangeFailed)
	}

	user, err := u.userRepo.UpsertSocialUser(ctx, entity.User{
		Email:       socialUser.Email,
		Name:        helper.TruncateString(socialUser.FullName(), 64, ""),
		CurrentStep: entity.CurrentStepInitial,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to upsert social user")
		return res, err
	}

	if user.Status == entity.UserStatusInactive {
		return res, domainError.New(code.UserStatusInactive)
	}

	if user.Status == entity.UserStatusBanned {
		return res, domainError.New(code.UserStatusBanned)
	}

	if user.CurrentStep == entity.CurrentStepCompleted {
		res.IsCompleted = true
	}

	rt, err := u.generateRefreshToken()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate refresh token")
		return res, err
	}

	sessionId, err := u.sessionRepo.Create(ctx, entity.Session{
		UserID:           user.ID,
		RefreshTokenHash: rt.HashRefreshToken,
		DeviceName:       helper.TruncateString(req.DeviceName, 128, ""),
		IPAddress:        req.IPAddress,
		ExpiredAt:        rt.ExpiresAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create session")
		return res, err
	}

	at, err := u.generateAccessToken(user.ID, sessionId)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate access token")
		return res, err
	}

	res.TokenResponse = TokenResponse{
		AccessToken:           at.AccessToken,
		AccessTokenExpiresAt:  at.ExpiresAt,
		RefreshToken:          rt.RefreshToken,
		RefreshTokenExpiresAt: rt.ExpiresAt,
	}

	return res, nil
}
