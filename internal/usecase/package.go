package usecase

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/auth"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/permission"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/role"
	storage "github.com/rizalarfiyan/be-plant-factory/internal/usecase/storage"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/tenant"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/user"
	userMe "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-me"
	usertenant "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-tenant"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(health.NewHealthUseCase),
	do.Lazy(auth.NewAuthUseCase),
	do.Lazy(userMe.NewUserMeUseCase),
	do.Lazy(user.NewUserUseCase),
	do.Lazy(role.NewRoleUseCase),
	do.Lazy(tenant.NewTenantUseCase),
	do.Lazy(usertenant.NewUserTenantUseCase),
	do.Lazy(permission.NewPermissionUseCase),
	do.Lazy(storage.NewStorageUseCase),
)
