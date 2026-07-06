package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u tenantUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TenantStatus) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	isUpdated, err := u.tenantRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to update tenant status")
		return err
	}

	if !isUpdated {
		return domainError.New(code.TenantHasSameStatus)
	}

	return nil
}
