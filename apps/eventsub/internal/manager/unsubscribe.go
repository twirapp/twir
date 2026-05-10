package manager

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *Manager) UnsubscribeChannel(ctx context.Context, channelID string) error {
	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		return err
	}

	c.logger.Info("unsubscribe twitch subscriptions: started", slog.String("channel_id", channelID))

	var cursor string
	var subscriptionIDsToRemove []string
	var scannedSubscriptions int
	for {
		existedSubsRes, err := twitchClient.GetEventSubSubscriptions(
			&helix.EventSubSubscriptionsParams{
				After: cursor,
			},
		)
		if err != nil {
			return err
		}

		if existedSubsRes == nil {
			return nil
		}

		if existedSubsRes.ErrorMessage != "" {
			return errors.New(existedSubsRes.ErrorMessage)
		}

		scannedSubscriptions += len(existedSubsRes.Data.EventSubSubscriptions)

		for _, sub := range existedSubsRes.Data.EventSubSubscriptions {
			if !shouldUnsubscribeChannelSubscription(channelID, sub) {
				continue
			}

			subscriptionIDsToRemove = append(subscriptionIDsToRemove, sub.ID)
		}

		cursor = existedSubsRes.Data.Pagination.Cursor
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
		res, err := twitchClient.RemoveEventSubSubscription(subscriptionID)
		if isSubscriptionNotFound(err, res) {
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

		if res != nil && res.ErrorMessage != "" {
			c.logger.Warn(
				"failed to remove subscription",
				slog.String("subscription_id", subscriptionID),
				slog.String("error", res.ErrorMessage),
			)
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

	return condition.BroadcasterUserID == channelID ||
		condition.UserID == channelID ||
		condition.ModeratorUserID == channelID ||
		condition.ToBroadcasterUserID == channelID ||
		condition.FromBroadcasterUserID == channelID
}

func isSubscriptionNotFound(err error, res *helix.RemoveEventSubSubscriptionParamsResponse) bool {
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "not found") {
		return true
	}

	if res == nil {
		return false
	}

	return res.StatusCode == 404 || strings.Contains(strings.ToLower(res.ErrorMessage), "not found")
}
