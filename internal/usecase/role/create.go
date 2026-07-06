package role

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (u roleUseCase) Create(ctx context.Context, role entity.Role) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	err := u.roleRepo.Create(ctx, role)
	if err != nil {
		log.Error().Err(err).Msg("failed to create role")
		return err
	}

	return nil
}
