package user

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u userUseCase) Update(ctx context.Context, user entity.User) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	isUpdated, err := u.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return err
	}

	if !isUpdated {
		return domainError.New(code.UserNotFound)
	}

	if err := u.userCache.Clear(ctx, user.ID); err != nil {
		log.Error().Err(err).Msg("failed to clear user cache")
		return err
	}

	return nil
}
