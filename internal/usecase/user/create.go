package user

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (u userUseCase) Create(ctx context.Context, user entity.User) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	err := u.userRepo.Create(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}

	return nil
}
