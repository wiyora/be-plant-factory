package storage

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (uc storageUseCase) CleanupTemporary(ctx context.Context) func() {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	return func() {
		if err := uc.storageRepo.CleanupTemporary(ctx); err != nil {
			log.Error().Err(err).Msg("failed to cleanup temporary files")
			return
		}
	}
}
