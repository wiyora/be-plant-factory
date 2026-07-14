package tenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
)

func (u tenantUseCase) Create(ctx context.Context, req entity.Tenant) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	diff := entity.StorageTypeTenantLogo.Diff("", req.Logo)
	if !diff.IsValid {
		return domainError.NewManualValidation("logo", "INVALID")
	}

	req.Logo = diff.Result
	err := u.tenantRepo.Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tenant")
		return err
	}

	if err := u.s3Repo.Diff(ctx, diff); err != nil {
		log.Error().Err(err).Msg("failed to process storage diff")
		return err
	}

	return nil
}
