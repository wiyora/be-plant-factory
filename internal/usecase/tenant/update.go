package tenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	domainError "github.com/rizalarfiyan/be-plant-factory/internal/domain/error"
	"github.com/rizalarfiyan/be-plant-factory/internal/infrastructure/logger"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/code"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/helper"
)

func (u tenantUseCase) Update(ctx context.Context, req entity.Tenant) error {
	log := logger.WithLayerCtx(ctx, logger.LayerUseCase)

	existing, err := u.tenantRepo.GetById(ctx, req.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get tenant by id")
		return err
	}

	if helper.IsEmptyStruct(existing) {
		return domainError.New(code.TenantNotFound)
	}

	diff := entity.StorageTypeTenantLogo.Diff(existing.Logo, req.Logo)
	if !diff.IsValid {
		return domainError.NewManualValidation("logo", "INVALID")
	}

	req.Logo = diff.Result
	isUpdated, err := u.tenantRepo.Update(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to update tenant")
		return err
	}

	if !isUpdated {
		return domainError.New(code.TenantNotFound)
	}

	if err := u.s3Repo.Diff(ctx, diff); err != nil {
		log.Error().Err(err).Msg("failed to process storage diff")
		return err
	}

	return nil
}
