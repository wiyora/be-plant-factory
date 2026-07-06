package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u tenantUseCase) Detail(ctx context.Context, id uuid.UUID) (entity.Tenant, error) {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	item, err := u.tenantRepo.GetById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get tenant by id")
		return entity.Tenant{}, err
	}

	if helper.IsEmptyStruct(item) {
		return entity.Tenant{}, domainError.New(code.NotFound)
	}

	return item.ToEntity(), nil
}
