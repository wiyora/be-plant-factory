package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u roleUseCase) SelectedDropdown(ctx context.Context, selectedIds []uuid.UUID) ([]entity.DropdownItem, error) {
	if len(selectedIds) == 0 {
		return nil, nil
	}

	items, err := u.roleRepo.SelectedDropdown(ctx, selectedIds)
	if err != nil {
		return nil, err
	}

	dropdowns := make([]entity.DropdownItem, len(items))
	for i, item := range items {
		dropdowns[i] = item.ToEntity()
	}

	return dropdowns, nil
}
