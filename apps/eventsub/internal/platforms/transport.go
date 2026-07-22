package platforms

import (
	"context"
	"errors"
	"fmt"

	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/libs/entities/platform"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type EventTransport interface {
	Platform() platform.Platform
	Capabilities() platform.Capabilities
	Subscribe(context.Context, channelplatformsmodel.ChannelPlatform) error
	Unsubscribe(context.Context, channelplatformsmodel.ChannelPlatform) error
	SetCallbackBaseURL(string)
}

func NewRegistry(transports ...EventTransport) *platformsregistry.Registry[EventTransport] {
	registry := platformsregistry.New[EventTransport]()
	for _, transport := range transports {
		registry.Register(transport)
	}

	return registry
}

func NewKickRegistry(kickSubManager *kick.SubscriptionManager) *platformsregistry.Registry[EventTransport] {
	return NewRegistry(kickSubManager)
}

func SubscribeAll(
	ctx context.Context,
	registry *platformsregistry.Registry[EventTransport],
	bindings []channelplatformsmodel.ChannelPlatform,
) error {
	var subscribeErrors []error
	for _, binding := range bindings {
		if !binding.Enabled {
			continue
		}

		transport, ok := registry.Get(binding.Platform)
		if !ok {
			continue
		}

		if err := transport.Subscribe(ctx, binding); err != nil {
			subscribeErrors = append(subscribeErrors, fmt.Errorf(
				"subscribe %q binding %q: %w",
				binding.Platform,
				binding.ID,
				err,
			))
		}
	}

	return errors.Join(subscribeErrors...)
}

func UnsubscribeAll(
	ctx context.Context,
	registry *platformsregistry.Registry[EventTransport],
	bindings []channelplatformsmodel.ChannelPlatform,
) error {
	var unsubscribeErrors []error
	for _, binding := range bindings {
		transport, ok := registry.Get(binding.Platform)
		if !ok {
			continue
		}

		if err := transport.Unsubscribe(ctx, binding); err != nil {
			unsubscribeErrors = append(unsubscribeErrors, fmt.Errorf(
				"unsubscribe %q binding %q: %w",
				binding.Platform,
				binding.ID,
				err,
			))
		}
	}

	return errors.Join(unsubscribeErrors...)
}
