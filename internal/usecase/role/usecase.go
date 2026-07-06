package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/samber/do/v2"
)

type RoleUseCase interface {
	List(ctx context.Context, filter entity.RoleFilter) ([]entity.Role, uint64, error)
	Detail(ctx context.Context, id uuid.UUID) (entity.Role, error)
	Create(ctx context.Context, role entity.Role) error
	Update(ctx context.Context, role entity.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllPermissions(ctx context.Context) []entity.PermissionModule
}

type roleUseCase struct {
	roleRepo postgres.RoleRepository `do:""`
}

func NewRoleUseCase(i do.Injector) (RoleUseCase, error) {
	return do.InvokeStruct[roleUseCase](i)
}
