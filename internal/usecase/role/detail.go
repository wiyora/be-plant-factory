package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u roleUseCase) Detail(ctx context.Context, id uuid.UUID) (entity.Role, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	item, err := u.roleRepo.GetById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get role by id")
		return entity.Role{}, err
	}

	if helper.IsEmptyStruct(item) {
		return entity.Role{}, domainError.New(code.RoleNotFound)
	}

	return item.ToEntity(), nil
}
