package resolvers

import (
	"context"
	"fmt"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	chatmessagesrepo "github.com/twirapp/twir/libs/repositories/chat_messages"
)

type apiKeyChannelIdentity struct {
	InternalChannelID string
	ChatTargets       []chatmessagesrepo.PlatformChannelIdentity
}

func resolveApiKeyChannelIdentity(
	ctx context.Context,
	deps Deps,
	apiKey string,
) (apiKeyChannelIdentity, error) {
	user, err := deps.UsersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return apiKeyChannelIdentity{}, fmt.Errorf("failed to get user: %w", err)
	}

	var channel channelsmodel.Channel
	switch user.Platform {
	case platformentity.PlatformKick:
		channel, err = deps.ChannelsRepository.GetByKickUserID(ctx, user.ID)
		if err != nil {
			return apiKeyChannelIdentity{}, fmt.Errorf("failed to get kick channel: %w", err)
		}
	default:
		channel, err = deps.ChannelsRepository.GetByTwitchUserID(ctx, user.ID)
		if err != nil {
			return apiKeyChannelIdentity{}, fmt.Errorf("failed to get twitch channel: %w", err)
		}
	}

	targets := make([]chatmessagesrepo.PlatformChannelIdentity, 0, 2)
	if channel.TwitchPlatformID != nil && *channel.TwitchPlatformID != "" {
		targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{
			Platform:          "twitch",
			PlatformChannelID: *channel.TwitchPlatformID,
		})
	}

	if channel.KickPlatformID != nil && *channel.KickPlatformID != "" {
		targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{
			Platform:          "kick",
			PlatformChannelID: *channel.KickPlatformID,
		})
	}

	if len(targets) == 0 {
		return apiKeyChannelIdentity{}, fmt.Errorf("no chat message targets found for api key")
	}

	return apiKeyChannelIdentity{
		InternalChannelID: channel.ID.String(),
		ChatTargets:       targets,
	}, nil
}
