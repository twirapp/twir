package platforms

import (
	"errors"
	"testing"

	"github.com/twirapp/twir/libs/entities/platform"
)

type testProvider struct {
	platform     platform.Platform
	capabilities platform.Capabilities
}

func (p testProvider) Platform() platform.Platform {
	return p.platform
}

func (p testProvider) Capabilities() platform.Capabilities {
	return p.capabilities
}

func TestRegistryGetsRegisteredProvider(t *testing.T) {
	t.Parallel()

	registry := New[testProvider]()
	registry.Register(testProvider{
		platform:     platform.PlatformTwitch,
		capabilities: platform.Capabilities{platform.CapabilityChatWrite},
	})

	provider, found := registry.Get(platform.PlatformTwitch)
	if !found {
		t.Fatal("expected registered Twitch provider")
	}
	if !provider.Capabilities().Supports(platform.CapabilityChatWrite) {
		t.Error("expected registered Twitch provider to support chat.write")
	}
}

func TestRegistryRequireReturnsUnsupportedCapability(t *testing.T) {
	t.Parallel()

	registry := New[testProvider]()
	registry.Register(testProvider{
		platform:     platform.PlatformKick,
		capabilities: platform.Capabilities{platform.CapabilityChatRead},
	})

	_, err := registry.Require(platform.PlatformKick, platform.CapabilityChatWrite)
	if err == nil {
		t.Fatal("expected chat.write to be rejected")
	}

	var unsupported platform.ErrUnsupportedCapability
	if !errors.As(err, &unsupported) {
		t.Fatalf("expected ErrUnsupportedCapability, got %v", err)
	}
	if unsupported.Platform != platform.PlatformKick || unsupported.Capability != platform.CapabilityChatWrite {
		t.Errorf("unexpected unsupported capability error: %#v", unsupported)
	}
}
