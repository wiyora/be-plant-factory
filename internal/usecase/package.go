package usecase

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/auth"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/permission"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/role"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/user"
	userMe "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-me"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(health.NewHealthUseCase),
	do.Lazy(auth.NewAuthUseCase),
	do.Lazy(userMe.NewUserMeUseCase),
	do.Lazy(user.NewUserUseCase),
	do.Lazy(role.NewRoleUseCase),
	do.Lazy(permission.NewPermissionUseCase),
)
