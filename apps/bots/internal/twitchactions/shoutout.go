package twitchactions

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

type ShoutOutInput struct {
	BroadcasterID string
	TargetID      string
}

func (c *TwitchActions) ShoutOut(ctx context.Context, input ShoutOutInput) error {
	channel, err := c.channelsByTwitchIDCache.Get(ctx, input.BroadcasterID)
	if err != nil {
		return fmt.Errorf("cannot get channel: %w", err)
	}
	if !channel.IsEnabled || !channel.IsBotMod || channel.IsTwitchBanned {
		return nil
	}
	if channel.TwitchUserID == nil {
		return fmt.Errorf("channel has no twitch user id for broadcaster %s", input.BroadcasterID)
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		*channel.TwitchPlatformID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return fmt.Errorf("cannot create broadcaster twitch client: %w", err)
	}

	resp, err := twitchClient.SendShoutout(
		&helix.SendShoutoutParams{
			FromBroadcasterID: input.BroadcasterID,
			ToBroadcasterID:   input.TargetID,
			ModeratorID:       input.BroadcasterID,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send shoutout: %w", err)
	}
	if resp.ErrorMessage != "" {
		return fmt.Errorf("cannot send shoutout: %s", resp.ErrorMessage)
	}

	return nil
}
