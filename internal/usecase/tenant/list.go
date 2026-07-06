package tenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u tenantUseCase) List(ctx context.Context, filter entity.TenantFilter) ([]entity.Tenant, uint64, error) {
	items, total, err := u.tenantRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	tenants := make([]entity.Tenant, len(items))
	for i, item := range items {
		tenants[i] = item.ToEntity()
	}

	return tenants, total, nil
}
