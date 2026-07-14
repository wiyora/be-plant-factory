package tenant

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u tenantUseCase) Dropdown(ctx context.Context, filter entity.DropdownFilter) ([]entity.DropdownItem, uint64, error) {
	items, total, err := u.tenantRepo.Dropdown(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	dropdowns := make([]entity.DropdownItem, len(items))
	for i, item := range items {
		dropdowns[i] = item.ToEntity()
	}

	return dropdowns, total, nil
}
