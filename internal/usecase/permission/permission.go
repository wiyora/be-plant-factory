package permission

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/shared/constant"
)

func (u permissionUseCase) All(ctx context.Context) []entity.PermissionModule {
	return constant.PermissionModules
}
