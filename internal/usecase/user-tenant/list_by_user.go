package usertenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u userTenantUseCase) ListByUser(ctx context.Context, filter entity.UserTenantFilter) ([]entity.UserTenantList, uint64, error) {
	items, total, err := u.userTenantRepo.ListByUser(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	tenants := make([]entity.UserTenantList, len(items))
	for i, item := range items {
		tenants[i] = item.ToEntity()
	}

	return tenants, total, nil
}
