package platforms

import "github.com/twirapp/twir/libs/entities/platform"

type Provider interface {
	Platform() platform.Platform
	Capabilities() platform.Capabilities
}

type Registry[T Provider] struct {
	providers map[platform.Platform]T
}

func New[T Provider]() *Registry[T] {
	return &Registry[T]{
		providers: make(map[platform.Platform]T),
	}
}

func (r *Registry[T]) Register(provider T) {
	r.providers[provider.Platform()] = provider
}

func (r *Registry[T]) Get(p platform.Platform) (T, bool) {
	provider, ok := r.providers[p]
	return provider, ok
}

func (r *Registry[T]) Require(p platform.Platform, capability platform.Capability) (T, error) {
	provider, ok := r.Get(p)
	if !ok || !provider.Capabilities().Supports(capability) {
		var zero T
		return zero, platform.ErrUnsupportedCapability{
			Platform:   p,
			Capability: capability,
		}
	}

	return provider, nil
}
