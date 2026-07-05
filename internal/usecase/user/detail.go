package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u userUseCase) Detail(ctx context.Context, id uuid.UUID) (entity.User, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	item, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by id")
		return entity.User{}, err
	}

	if helper.IsEmptyStruct(item) {
		return entity.User{}, domainError.New(code.UserNotFound)
	}

	return item.ToEntity(), nil
}
