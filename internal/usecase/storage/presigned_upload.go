package storage

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (uc storageUseCase) GeneratePresignedUpload(ctx context.Context, req entity.StorageGeneratePresignedUpload) (entity.StoragePresignedUpload, error) {
	result, err := uc.storageRepo.GeneratePresignedUpload(ctx, req)
	if err != nil {
		return entity.StoragePresignedUpload{}, err
	}

	return result.ToEntity(), nil
}
