package storage

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/storage"
	"github.com/samber/do/v2"
)

type StorageUseCase interface {
	GeneratePresignedUpload(ctx context.Context, req entity.StorageGeneratePresignedUpload) (entity.StoragePresignedUpload, error)
	CleanupTemporary(ctx context.Context) func()
}

type storageUseCase struct {
	storageRepo storage.S3Repository `do:""`
}

func NewStorageUseCase(i do.Injector) (StorageUseCase, error) {
	return do.InvokeStruct[storageUseCase](i)
}
