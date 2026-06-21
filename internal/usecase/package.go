package usecase

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/usecase/health"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(health.NewHealthUseCase),
)
