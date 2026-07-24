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
	twitchBinding, botConfig, found, err := channel.TwitchBinding()
	if err != nil {
		return fmt.Errorf("cannot parse Twitch bot config: %w", err)
	}
	if !found || !twitchBinding.Enabled || !botConfig.IsBotMod || botConfig.IsTwitchBanned ||
		twitchBinding.PlatformChannelID == "" {
		return nil
	}
	if twitchBinding.PlatformChannelID != input.BroadcasterID {
		return fmt.Errorf("Twitch binding channel id does not match broadcaster %s", input.BroadcasterID)
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchBinding.UserID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return fmt.Errorf("cannot create broadcaster twitch client: %w", err)
	}

	resp, err := twitchClient.SendShoutout(
		&helix.SendShoutoutParams{
			FromBroadcasterID: twitchBinding.PlatformChannelID,
			ToBroadcasterID:   input.TargetID,
			ModeratorID:       twitchBinding.PlatformChannelID,
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
