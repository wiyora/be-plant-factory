package social

import (
	"errors"

	"github.com/rizalarfiyan/be-plant-factory/internal/domain/entity"
)

var ErrProviderNotFound = errors.New("social auth provider not found")

type Manager interface {
	GetProvider(name entity.Provider) (SocialAuthProvider, error)
	Register(name entity.Provider, provider SocialAuthProvider)
}

type manager struct {
	providers map[entity.Provider]SocialAuthProvider
}

func NewManager() Manager {
	return &manager{
		providers: make(map[entity.Provider]SocialAuthProvider),
	}
}

func (m *manager) Register(name entity.Provider, provider SocialAuthProvider) {
	m.providers[name] = provider
}

func (m *manager) GetProvider(name entity.Provider) (SocialAuthProvider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, ErrProviderNotFound
	}

	return provider, nil
}
