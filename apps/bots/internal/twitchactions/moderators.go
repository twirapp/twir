package twitchactions

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/bots/internal/channelbinding"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *TwitchActions) AddModerator(ctx context.Context, broadcasterID, userID string) error {
	channel, err := c.channelsByTwitchIDCache.Get(ctx, broadcasterID)
	if err != nil {
		return fmt.Errorf("cannot get channel by twitch id: %w", err)
	}
	twitchBinding, found := channelbinding.Find(channel, platform.PlatformTwitch)
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" {
		return fmt.Errorf("channel has no enabled Twitch binding for broadcaster %s", broadcasterID)
	}
	if twitchBinding.PlatformChannelID != broadcasterID {
		return fmt.Errorf("Twitch binding channel id does not match broadcaster %s", broadcasterID)
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchBinding.UserID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
	}

	resp, err := twitchClient.AddChannelModerator(
		&helix.AddChannelModeratorParams{
			BroadcasterID: twitchBinding.PlatformChannelID,
			UserID:        userID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add moderator: %w", err)
	}
	if resp.ErrorMessage != "" {
		return fmt.Errorf("failed to add moderator: %s", resp.ErrorMessage)
	}

	return nil
}

func (c *TwitchActions) RemoveModerator(ctx context.Context, broadcasterID, userID string) error {
	channel, err := c.channelsByTwitchIDCache.Get(ctx, broadcasterID)
	if err != nil {
		return fmt.Errorf("cannot get channel by twitch id: %w", err)
	}
	twitchBinding, found := channelbinding.Find(channel, platform.PlatformTwitch)
	if !found || !twitchBinding.Enabled || twitchBinding.PlatformChannelID == "" {
		return fmt.Errorf("channel has no enabled Twitch binding for broadcaster %s", broadcasterID)
	}
	if twitchBinding.PlatformChannelID != broadcasterID {
		return fmt.Errorf("Twitch binding channel id does not match broadcaster %s", broadcasterID)
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		twitchBinding.UserID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
	}

	resp, err := twitchClient.RemoveChannelModerator(
		&helix.RemoveChannelModeratorParams{
			BroadcasterID: twitchBinding.PlatformChannelID,
			UserID:        userID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to remove moderator: %w", err)
	}
	if resp.ErrorMessage != "" {
		return fmt.Errorf("failed to remove moderator: %s", resp.ErrorMessage)
	}

	return nil
}
