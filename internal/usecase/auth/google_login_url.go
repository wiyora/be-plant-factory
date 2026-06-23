package auth

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/redis"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u authUseCase) GoogleLoginURL(ctx context.Context, req GoogleRequest) (string, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	state, err := helper.GenerateRandomToken(32)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate random token for state")
		return "", err
	}

	statePayload := redis.StatePayload{
		State:       state,
		BaseURL:     req.BaseUrl,
		RedirectUri: req.RedirectUri,
	}
	if err := u.stateRepo.Set(ctx, entity.ProviderGoogle, statePayload); err != nil {
		log.Error().Err(err).Msg("failed to set state")
		return "", err
	}

	provider, err := u.authManager.GetProvider(entity.ProviderGoogle)
	if err != nil {
		log.Error().Err(err).Msg("failed to get Google provider")
		return "", err
	}

	return provider.GetLoginURL(state), nil
}
