package manager

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

func (c *Manager) unsubscribeFromTopic(ctx context.Context, channelID, topic string) error {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	existedSubsRes, _ := twitchClient.GetEventSubSubscriptions(
		&helix.EventSubSubscriptionsParams{
			Type:   topic,
			UserID: channelID,
		},
	)

	if len(existedSubsRes.Data.EventSubSubscriptions) > 0 {
		for _, sub := range existedSubsRes.Data.EventSubSubscriptions {
			res, err := twitchClient.RemoveEventSubSubscription(sub.ID)
			if err != nil {
				return fmt.Errorf("failed to remove subscription: %w", err)
			}
			if res.ErrorMessage != "" {
				return fmt.Errorf("failed to remove subscription: %s", res.ErrorMessage)
			}
		}
	}

	return nil
}
