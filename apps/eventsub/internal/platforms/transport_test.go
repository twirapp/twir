package platforms

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	channelplatformsmodel "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

type recordingTransport struct {
	platform     platform.Platform
	capabilities platform.Capabilities
	subscribed   []channelplatformsmodel.ChannelPlatform
	unsubscribed []channelplatformsmodel.ChannelPlatform
	subscribeErr map[uuid.UUID]error
}

var _ EventTransport = (*recordingTransport)(nil)

func (t *recordingTransport) Platform() platform.Platform {
	return t.platform
}

func (t *recordingTransport) Capabilities() platform.Capabilities {
	return t.capabilities
}

func (t *recordingTransport) Subscribe(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
) error {
	t.subscribed = append(t.subscribed, binding)
	return t.subscribeErr[binding.ID]
}

func (t *recordingTransport) Unsubscribe(
	_ context.Context,
	binding channelplatformsmodel.ChannelPlatform,
) error {
	t.unsubscribed = append(t.unsubscribed, binding)
	return nil
}

func (*recordingTransport) SetCallbackBaseURL(string) {}

func TestSubscribeAllRoutesBindingsToTheirMatchingTransport(t *testing.T) {
	twitchBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformTwitch,
		Enabled:  true,
	}
	kickBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformKick,
		Enabled:  true,
	}
	twitch := &recordingTransport{platform: platform.PlatformTwitch}
	kick := &recordingTransport{platform: platform.PlatformKick}

	err := SubscribeAll(
		context.Background(),
		NewRegistry(twitch, kick),
		[]channelplatformsmodel.ChannelPlatform{kickBinding, twitchBinding},
	)
	if err != nil {
		t.Fatalf("SubscribeAll returned error: %v", err)
	}
	if len(twitch.subscribed) != 1 || twitch.subscribed[0].ID != twitchBinding.ID {
		t.Errorf("Twitch subscriptions = %#v, want [%s]", twitch.subscribed, twitchBinding.ID)
	}
	if len(kick.subscribed) != 1 || kick.subscribed[0].ID != kickBinding.ID {
		t.Errorf("Kick subscriptions = %#v, want [%s]", kick.subscribed, kickBinding.ID)
	}
}

func TestSubscribeAllSkipsDisabledBindings(t *testing.T) {
	enabledBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformKick,
		Enabled:  true,
	}
	disabledBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformKick,
		Enabled:  false,
	}
	transport := &recordingTransport{platform: platform.PlatformKick}

	err := SubscribeAll(
		context.Background(),
		NewRegistry(transport),
		[]channelplatformsmodel.ChannelPlatform{disabledBinding, enabledBinding},
	)
	if err != nil {
		t.Fatalf("SubscribeAll returned error: %v", err)
	}
	if len(transport.subscribed) != 1 || transport.subscribed[0].ID != enabledBinding.ID {
		t.Errorf("transport subscriptions = %#v, want [%s]", transport.subscribed, enabledBinding.ID)
	}
}

func TestSubscribeAllContinuesAfterBindingFailure(t *testing.T) {
	failedBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformKick,
		Enabled:  true,
	}
	followingBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformKick,
		Enabled:  true,
	}
	wantErr := errors.New("subscribe failed")
	transport := &recordingTransport{
		platform: platform.PlatformKick,
		subscribeErr: map[uuid.UUID]error{
			failedBinding.ID: wantErr,
		},
	}

	err := SubscribeAll(
		context.Background(),
		NewRegistry(transport),
		[]channelplatformsmodel.ChannelPlatform{failedBinding, followingBinding},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("SubscribeAll error = %v, want %v", err, wantErr)
	}
	if len(transport.subscribed) != 2 || transport.subscribed[1].ID != followingBinding.ID {
		t.Errorf("transport subscriptions = %#v, want failed and following bindings", transport.subscribed)
	}
}
