package role

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u roleUseCase) Update(ctx context.Context, role entity.Role) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	isUpdated, err := u.roleRepo.Update(ctx, role)
	if err != nil {
		log.Error().Err(err).Msg("failed to update role")
		return err
	}

	if !isUpdated {
		return domainError.New(code.RoleNotFound)
	}

	return nil
}
