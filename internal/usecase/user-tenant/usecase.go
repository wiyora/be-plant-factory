package usertenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/samber/do/v2"
)

type UserTenantUseCase interface {
	ListByTenant(ctx context.Context, filter entity.TenantUserFilter) ([]entity.TenantUserList, uint64, error)
	ListByUser(ctx context.Context, filter entity.UserTenantFilter) ([]entity.UserTenantList, uint64, error)
	AssignRole(ctx context.Context, tenantID, userID, roleID uuid.UUID) error
	RemoveUser(ctx context.Context, tenantID, userID uuid.UUID) error
}

type userTenantUseCase struct {
	userTenantRepo postgres.UserTenantRepository `do:""`
}

func NewUserTenantUseCase(i do.Injector) (UserTenantUseCase, error) {
	return do.InvokeStruct[userTenantUseCase](i)
}
