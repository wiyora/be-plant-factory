package role

import (
	"context"

	"github.com/google/uuid"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u roleUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	isDeleted, err := u.roleRepo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete role")
		return err
	}

	if !isDeleted {
		return domainError.New(code.RoleNotFound)
	}

	return nil
}
