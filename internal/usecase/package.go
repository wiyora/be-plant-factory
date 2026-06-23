package usecase

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/auth"
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
	userMe "github.com/rizalarfiyan/be-plant-factory/internal/usecase/user-me"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(health.NewHealthUseCase),
	do.Lazy(auth.NewAuthUseCase),
	do.Lazy(userMe.NewUserMeUseCase),
)
