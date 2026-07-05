package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u userUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	isUpdated, err := u.userRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to update user status")
		return err
	}

	if !isUpdated {
		return domainError.New(code.UserHasSameStatus)
	}

	if err := u.userCache.Clear(ctx, id); err != nil {
		log.Error().Err(err).Msg("failed to clear user cache")
		return err
	}

	return nil
}
