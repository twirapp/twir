package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *Service) GetChannelChatBadges(ctx context.Context, channelID string) (
	[]helix.ChatBadge,
	error,
) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelsRepository.GetByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() || !channel.TwitchConnected() {
		return nil, nil
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: *channel.TwitchPlatformID,
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf(
			"cannot get channel badges: %v %s",
			resp.StatusCode,
			resp.ErrorMessage,
		)
	}

	return resp.Data.Badges, nil
}

func (c *Service) GetGlobalChatBadges(ctx context.Context) (
	[]helix.ChatBadge,
	error,
) {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.GetGlobalChatBadges()
	if err != nil {
		return nil, err
	}
	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf(
			"cannot get global badges: %v %s",
			resp.StatusCode,
			resp.ErrorMessage,
		)
	}

	return resp.Data.Badges, nil
}
