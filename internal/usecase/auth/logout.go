package auth

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (u authUseCase) Logout(ctx context.Context, req LogoutRequest) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	if err := u.sessionRepo.Delete(ctx, req.UserID, req.SessionID); err != nil {
		log.Error().Err(err).Msg("failed to delete session")
		return err
	}

	if err := u.tokenRepo.Blacklist(ctx, req.AccessTokenID, req.ExpiredAt); err != nil {
		log.Error().Err(err).Msg("failed to blacklist access token")
		return err
	}

	return nil
}
