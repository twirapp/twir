package platforms

import (
	"context"
	"errors"
	"fmt"

	kickchat "github.com/twirapp/twir/apps/bots/internal/kick"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/platform"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type ChatOptions struct {
	IsAnnounce        bool
	SkipToxicityCheck bool
	SkipRateLimits    bool
	AnnounceColor     bots.AnnounceColor
}

type ChatAdapter interface {
	platformsregistry.Provider
	SendMessage(context.Context, channelplatformsmodel.ChannelPlatform, string, string, ChatOptions) error
}

func newRegistry(adapters ...ChatAdapter) *platformsregistry.Registry[ChatAdapter] {
	registry := platformsregistry.New[ChatAdapter]()
	for _, adapter := range adapters {
		registry.Register(adapter)
	}

	return registry
}

func NewChatRegistry(
	twitchActions *twitchactions.TwitchActions,
	kickClient *kickchat.ChatClient,
) *platformsregistry.Registry[ChatAdapter] {
	return newRegistry(
		NewTwitchChatAdapter(twitchActions),
		NewKickChatAdapter(kickClient),
	)
}

func Dispatch(
	ctx context.Context,
	registry *platformsregistry.Registry[ChatAdapter],
	bindings []channelplatformsmodel.ChannelPlatform,
	requestedPlatforms []platform.Platform,
	message string,
	replyID string,
	options ChatOptions,
) error {
	var dispatchErrors []error
	if len(requestedPlatforms) > 0 {
		for _, requestedPlatform := range requestedPlatforms {
			adapter, err := registry.Require(requestedPlatform, platform.CapabilityChatWrite)
			if err != nil {
				dispatchErrors = append(dispatchErrors, fmt.Errorf(
					"dispatch chat message to platform %q: %w",
					requestedPlatform,
					err,
				))
				continue
			}

			for _, binding := range bindings {
				if binding.Platform != requestedPlatform {
					continue
				}
				if binding.Enabled {
					if err := adapter.SendMessage(ctx, binding, message, replyID, options); err != nil {
						dispatchErrors = append(dispatchErrors, fmt.Errorf(
							"send chat message to %q binding %q: %w",
							binding.Platform,
							binding.PlatformChannelID,
							err,
						))
					}
				}

				break
			}
		}

		return errors.Join(dispatchErrors...)
	}

	for _, binding := range bindings {
		if !binding.Enabled {
			continue
		}

		adapter, err := registry.Require(binding.Platform, platform.CapabilityChatWrite)
		if err != nil {
			continue
		}

		if err := adapter.SendMessage(ctx, binding, message, replyID, options); err != nil {
			dispatchErrors = append(dispatchErrors, fmt.Errorf(
				"send chat message to %q binding %q: %w",
				binding.Platform,
				binding.PlatformChannelID,
				err,
			))
		}
	}

	return errors.Join(dispatchErrors...)
}

type twitchMessageSender interface {
	SendMessage(context.Context, channelplatformsmodel.ChannelPlatform, twitchactions.SendMessageOpts) error
}

type twitchChatAdapter struct {
	sender twitchMessageSender
}

func NewTwitchChatAdapter(sender twitchMessageSender) ChatAdapter {
	return twitchChatAdapter{sender: sender}
}

func (a twitchChatAdapter) Platform() platform.Platform {
	return platform.PlatformTwitch
}

func (a twitchChatAdapter) Capabilities() platform.Capabilities {
	return platform.Capabilities{
		platform.CapabilityChatWrite,
		platform.CapabilityChatReply,
	}
}

func (a twitchChatAdapter) SendMessage(
	ctx context.Context,
	binding channelplatformsmodel.ChannelPlatform,
	message string,
	replyID string,
	options ChatOptions,
) error {
	return a.sender.SendMessage(
		ctx,
		binding,
		twitchactions.SendMessageOpts{
			Message:              message,
			ReplyParentMessageID: replyID,
			IsAnnounce:           options.IsAnnounce,
			SkipToxicityCheck:    options.SkipToxicityCheck,
			SkipRateLimits:       options.SkipRateLimits,
			AnnounceColor:        options.AnnounceColor,
		},
	)
}

type kickMessageSender interface {
	SendMessage(context.Context, channelplatformsmodel.ChannelPlatform, string, string) error
}

type kickChatAdapter struct {
	sender kickMessageSender
}

func NewKickChatAdapter(sender kickMessageSender) ChatAdapter {
	return kickChatAdapter{sender: sender}
}

func (a kickChatAdapter) Platform() platform.Platform {
	return platform.PlatformKick
}

func (a kickChatAdapter) Capabilities() platform.Capabilities {
	return platform.Capabilities{
		platform.CapabilityChatWrite,
		platform.CapabilityChatReply,
	}
}

func (a kickChatAdapter) SendMessage(
	ctx context.Context,
	binding channelplatformsmodel.ChannelPlatform,
	message string,
	replyID string,
	_ ChatOptions,
) error {
	return a.sender.SendMessage(ctx, binding, message, replyID)
}
