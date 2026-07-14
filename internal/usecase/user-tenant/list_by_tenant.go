package usertenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u userTenantUseCase) ListByTenant(ctx context.Context, filter entity.TenantUserFilter) ([]entity.TenantUserList, uint64, error) {
	items, total, err := u.userTenantRepo.ListByTenant(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	users := make([]entity.TenantUserList, len(items))
	for i, item := range items {
		users[i] = item.ToEntity()
	}

	return users, total, nil
}
