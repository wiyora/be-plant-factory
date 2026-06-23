package social

import (
	"github.com/rizalarfiyan/be-plant-factory/internal/config"
	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
	"github.com/samber/do/v2"
)

func New(i do.Injector) (Manager, error) {
	cfg := do.MustInvoke[*config.Config](i)

	manager := NewManager()

	googleProvider := NewGoogleProvider(cfg)
	manager.Register(entity.ProviderGoogle, googleProvider)

	return manager, nil
}
