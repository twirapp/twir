package twitchactions

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix"
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

	twitchClient, err := c.createUserClient(ctx, twitchBinding.UserID)
	if err != nil {
		return fmt.Errorf("cannot create broadcaster twitch client: %w", err)
	}

	_, err = twitchClient.Chat.SendShoutout(
		ctx,
		helix.SendShoutoutRequest{
			FromBroadcasterID: twitchBinding.PlatformChannelID,
			ToBroadcasterID:   input.TargetID,
			ModeratorID:       twitchBinding.PlatformChannelID,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send shoutout: %w", err)
	}

	return nil
}
