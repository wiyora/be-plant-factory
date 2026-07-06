package permission

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/samber/do/v2"
)

type PermissionUseCase interface {
	All(ctx context.Context) []entity.PermissionModule
}

type permissionUseCase struct {
}

func NewPermissionUseCase(i do.Injector) (PermissionUseCase, error) {
	return do.InvokeStruct[permissionUseCase](i)
}
