package userMe

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
)

func (u userMeUseCase) GettingStarted(ctx context.Context, req entity.UserMeGettingStarted) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	if req.CurrentStep != entity.CurrentStepInitial {
		log.Warn().Str("current_step", req.CurrentStep.String()).Msg("invalid getting started step")
		return domainError.New(code.InvalidCurrentStep, domainError.WithParams(map[string]any{
			"current_step": req.CurrentStep.String(),
		}))
	}

	if err := u.userRepo.UpdateGettingStarted(ctx, req); err != nil {
		log.Error().Err(err).Msg("failed to update getting started")
		return err
	}

	if err := u.userCache.Clear(ctx, req.UserID); err != nil {
		log.Error().Err(err).Msg("failed to clear user cache")
		return err
	}

	return nil
}
