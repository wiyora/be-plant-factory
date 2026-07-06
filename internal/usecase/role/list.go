package role

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u roleUseCase) List(ctx context.Context, filter entity.RoleFilter) ([]entity.Role, uint64, error) {
	items, total, err := u.roleRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	roles := make([]entity.Role, len(items))
	for i, item := range items {
		roles[i] = item.ToEntity()
	}

	return roles, total, nil
}
