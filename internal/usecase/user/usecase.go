package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/cache"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/samber/do/v2"
)

type UserUseCase interface {
	List(ctx context.Context, filter entity.UserFilter) ([]entity.User, uint64, error)
	Detail(ctx context.Context, id uuid.UUID) (entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.UserStatus) error
	Dropdown(ctx context.Context, filter entity.DropdownFilter) ([]entity.DropdownItem, uint64, error)
	SelectedDropdown(ctx context.Context, selectedIds []uuid.UUID) ([]entity.DropdownItem, error)
}

type userUseCase struct {
	userRepo  postgres.UserRepository `do:""`
	userCache cache.UserRepository    `do:""`
}

func NewUserUseCase(i do.Injector) (UserUseCase, error) {
	return do.InvokeStruct[userUseCase](i)
}
