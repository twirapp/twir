package manager

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/libs/logger"
)

func (c *Manager) UnsubscribeChannel(ctx context.Context, channelID string) error {
	twitchClient, err := c.newAppTwitchClient(ctx)
	if err != nil {
		return err
	}

	c.logger.Info("unsubscribe twitch subscriptions: started", slog.String("channel_id", channelID))

	var cursor string
	var subscriptionIDsToRemove []string
	var scannedSubscriptions int
	for {
		existedSubsRes, err := twitchClient.EventSub.GetEventSubSubscriptions(
			ctx,
			helix.GetEventSubSubscriptionsRequest{
				After:           cursor,
				TransportMethod: helix.EventSubTransportConduit,
			},
		)
		if err != nil {
			return err
		}

		scannedSubscriptions += len(existedSubsRes.Data.Subscriptions)

		for _, sub := range existedSubsRes.Data.Subscriptions {
			if !shouldUnsubscribeChannelSubscription(channelID, sub) {
				continue
			}

			subscriptionIDsToRemove = append(subscriptionIDsToRemove, sub.ID)
		}

		cursor = existedSubsRes.Pagination.Cursor()
		if cursor == "" {
			break
		}
	}

	c.logger.Info(
		"unsubscribe twitch subscriptions: matched subscriptions",
		slog.String("channel_id", channelID),
		slog.Int("scanned_count", scannedSubscriptions),
		slog.Int("matched_count", len(subscriptionIDsToRemove)),
	)

	removedCount := 0
	notFoundCount := 0

	for _, subscriptionID := range subscriptionIDsToRemove {
		_, err := twitchClient.EventSub.DeleteEventSubSubscription(ctx, helix.DeleteEventSubSubscriptionRequest{
			ID:              subscriptionID,
			TransportMethod: helix.EventSubTransportConduit,
		})
		if isSubscriptionNotFound(err) {
			notFoundCount++
			c.logger.Info(
				"unsubscribe twitch subscriptions: subscription already absent",
				slog.String("channel_id", channelID),
				slog.String("subscription_id", subscriptionID),
			)
			continue
		}

		if err != nil {
			c.logger.Warn("failed to remove subscription", logger.Error(err), slog.String("subscription_id", subscriptionID))
			continue
		}

		removedCount++
	}

	c.logger.Info(
		"unsubscribe twitch subscriptions: finished",
		slog.String("channel_id", channelID),
		slog.Int("matched_count", len(subscriptionIDsToRemove)),
		slog.Int("removed_count", removedCount),
		slog.Int("already_absent_count", notFoundCount),
	)

	return nil
}

func shouldUnsubscribeChannelSubscription(channelID string, sub helix.EventSubSubscription) bool {
	condition := sub.Condition

	return condition["broadcaster_user_id"] == channelID ||
		condition["user_id"] == channelID ||
		condition["moderator_user_id"] == channelID ||
		condition["to_broadcaster_user_id"] == channelID ||
		condition["from_broadcaster_user_id"] == channelID
}

func isSubscriptionNotFound(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *helix.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode() == http.StatusNotFound {
		return true
	}

	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
