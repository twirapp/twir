package platform

import (
	"fmt"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type Registry struct {
	providers map[platformentity.Platform]PlatformProvider
}

func NewRegistry(providers []PlatformProvider) *Registry {
	registry := &Registry{providers: make(map[platformentity.Platform]PlatformProvider, len(providers))}
	for _, provider := range providers {
		if provider == nil {
			continue
		}

		registry.providers[platformentity.Platform(provider.Name())] = provider
	}

	return registry
}

func NewFeatureGatedRegistry(
	vkVideoEnabled bool,
	providers []PlatformProvider,
	newVKVideoProvider func() (PlatformProvider, error),
) (*Registry, error) {
	if !vkVideoEnabled {
		return NewRegistry(providers), nil
	}

	vkVideoProvider, err := newVKVideoProvider()
	if err != nil {
		return nil, fmt.Errorf("create VK Video provider: %w", err)
	}

	return NewRegistry(append(providers, vkVideoProvider)), nil
}

func (r *Registry) Get(platform platformentity.Platform) (PlatformProvider, bool) {
	provider, ok := r.providers[platform]
	return provider, ok
}
