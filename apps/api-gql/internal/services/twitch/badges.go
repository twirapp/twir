package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func (c *Service) GetChannelChatBadges(ctx context.Context, channelID string) (
	[]helix.ChatBadgeSet,
	error,
) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return nil, nil
	}

	twitchBinding, found := channel.Binding(platformentity.PlatformTwitch)
	if !found || twitchBinding.UserID == uuid.Nil {
		return nil, nil
	}

	twitchClient, err := c.createAppClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.Chat.GetChannelChatBadges(
		ctx, helix.GetChannelChatBadgesRequest{
			BroadcasterID: twitchBinding.PlatformChannelID,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Service) GetGlobalChatBadges(ctx context.Context) (
	[]helix.ChatBadgeSet,
	error,
) {
	twitchClient, err := c.createAppClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.Chat.GetGlobalChatBadges(ctx, helix.GetGlobalChatBadgesRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
