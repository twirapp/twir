package platform

import (
	"context"
	"errors"
	"testing"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestRegistryGetsRegisteredProvider(t *testing.T) {
	twitch := &fakeProvider{platform: platformentity.PlatformTwitch}
	registry := NewRegistry([]PlatformProvider{twitch})

	provider, ok := registry.Get(platformentity.PlatformTwitch)
	if !ok || provider != twitch {
		t.Fatalf("registry did not return the registered Twitch provider: %#v, %t", provider, ok)
	}
}

func TestRegistryDoesNotContainVKUntilItIsRegistered(t *testing.T) {
	registry := NewRegistry([]PlatformProvider{
		&fakeProvider{platform: platformentity.PlatformTwitch},
		&fakeProvider{platform: platformentity.PlatformKick},
	})

	if provider, ok := registry.Get(platformentity.PlatformVKVideoLive); ok || provider != nil {
		t.Fatalf("VK Video Live must not be registered while the feature is disabled: %#v", provider)
	}
}

func TestNewFeatureGatedRegistryDoesNotCreateVKProviderWhileDisabled(t *testing.T) {
	created := 0
	registry, err := NewFeatureGatedRegistry(
		false,
		[]PlatformProvider{&fakeProvider{platform: platformentity.PlatformTwitch}},
		func() (PlatformProvider, error) {
			created++
			return &fakeProvider{platform: platformentity.PlatformVKVideoLive}, nil
		},
	)
	if err != nil {
		t.Fatalf("create disabled registry: %v", err)
	}
	if created != 0 {
		t.Fatalf("VK provider factory called %d times while disabled", created)
	}
	if _, ok := registry.Get(platformentity.PlatformVKVideoLive); ok {
		t.Fatal("VK Video Live must not be registered while disabled")
	}
}

func TestNewFeatureGatedRegistryRegistersVKProviderWhenEnabled(t *testing.T) {
	vkProvider := &fakeProvider{platform: platformentity.PlatformVKVideoLive}
	registry, err := NewFeatureGatedRegistry(
		true,
		[]PlatformProvider{&fakeProvider{platform: platformentity.PlatformTwitch}},
		func() (PlatformProvider, error) {
			return vkProvider, nil
		},
	)
	if err != nil {
		t.Fatalf("create enabled registry: %v", err)
	}

	provider, ok := registry.Get(platformentity.PlatformVKVideoLive)
	if !ok || provider != vkProvider {
		t.Fatalf("VK Video Live was not registered: %#v, %t", provider, ok)
	}
}

func TestNewFeatureGatedRegistryReturnsVKFactoryError(t *testing.T) {
	wantErr := errors.New("invalid VK configuration")
	_, err := NewFeatureGatedRegistry(
		true,
		nil,
		func() (PlatformProvider, error) {
			return nil, wantErr
		},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected VK provider factory error %v, got %v", wantErr, err)
	}
}

type fakeProvider struct {
	platform platformentity.Platform
}

func (f *fakeProvider) Platform() platformentity.Platform {
	return f.platform
}

func (f *fakeProvider) GetAuthURL(string, string) string {
	return ""
}

func (f *fakeProvider) ExchangeCode(context.Context, ExchangeCodeInput) (*PlatformTokens, error) {
	return nil, nil
}

func (f *fakeProvider) RefreshToken(context.Context, RefreshTokenInput) (*PlatformTokens, error) {
	return nil, nil
}

func (f *fakeProvider) GetUser(context.Context, string) (*PlatformUser, error) {
	return nil, nil
}
