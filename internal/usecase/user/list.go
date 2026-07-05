package user

import (
	"context"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

func (u userUseCase) List(ctx context.Context, filter entity.UserFilter) ([]entity.User, uint64, error) {
	items, total, err := u.userRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, len(items))
	for i, item := range items {
		users[i] = item.ToEntity()
	}

	return users, total, nil
}
