package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	apiChannelbinding "github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return nil, nil
	}

	twitchBinding, found := apiChannelbinding.Find(channel, platformentity.PlatformTwitch)
	if !found || twitchBinding.UserID == uuid.Nil || twitchBinding.PlatformChannelID == "" {
		return nil, nil
	}

	twitchClient, err := c.createAppClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := twitchClient.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: twitchBinding.PlatformChannelID,
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
