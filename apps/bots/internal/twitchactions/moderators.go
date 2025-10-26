package twitchactions

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *TwitchActions) AddModerator(ctx context.Context, broadcasterID, userID string) error {
	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		broadcasterID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
	}

	resp, err := twitchClient.AddChannelModerator(
		&helix.AddChannelModeratorParams{
			BroadcasterID: broadcasterID,
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
	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		broadcasterID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
	}

	resp, err := twitchClient.RemoveChannelModerator(
		&helix.RemoveChannelModeratorParams{
			BroadcasterID: broadcasterID,
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
