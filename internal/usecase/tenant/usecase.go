package tenant

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/postgres"
	"github.com/rizalarfiyan/be-plant-factory/internal/repository/storage"
	"github.com/samber/do/v2"
)

type TenantUseCase interface {
	List(ctx context.Context, filter entity.TenantFilter) ([]entity.Tenant, uint64, error)
	Detail(ctx context.Context, id uuid.UUID) (entity.Tenant, error)
	Create(ctx context.Context, req entity.Tenant) error
	Update(ctx context.Context, req entity.Tenant) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.TenantStatus) error
	Dropdown(ctx context.Context, filter entity.DropdownFilter) ([]entity.DropdownItem, uint64, error)
}

type tenantUseCase struct {
	tenantRepo postgres.TenantRepository `do:""`
	s3Repo     storage.S3Repository      `do:""`
}

func NewTenantUseCase(i do.Injector) (TenantUseCase, error) {
	return do.InvokeStruct[tenantUseCase](i)
}
