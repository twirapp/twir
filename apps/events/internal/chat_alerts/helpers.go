package chat_alerts

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) takeCountedSample(
	target int,
	messages []model.ChatAlertsCountedMessage,
) string {
	if len(messages) == 0 {
		return ""
	}

	slices.SortFunc(
		messages, func(a, b model.ChatAlertsCountedMessage) int {
			return cmp.Compare(a.Count, b.Count)
		},
	)

	var lastMatch model.ChatAlertsCountedMessage
	for _, m := range messages {
		if m.Count <= target {
			lastMatch = m
		}
	}

	groupedMatched := lo.GroupBy(
		messages,
		func(m model.ChatAlertsCountedMessage) int {
			return m.Count
		},
	)
	sample := lo.Sample(groupedMatched[lastMatch.Count]).Text

	return sample
}

func (c *ChatAlerts) buildRedisCooldownKey(channelId, eventName string) string {
	return fmt.Sprintf(
		"channels:%s:chat_alerts_events:cooldowns:%s",
		channelId,
		eventName,
	)
}

func (c *ChatAlerts) isOnCooldown(ctx context.Context, channelId, eventName string) (
	bool,
	error,
) {
	exists, err := c.redis.Exists(
		ctx,
		c.buildRedisCooldownKey(channelId, eventName),
	).Result()
	return exists == 1, err
}

func (c *ChatAlerts) SetCooldown(
	ctx context.Context,
	channelId, eventName string,
	seconds int,
) error {
	if seconds == 0 {
		return nil
	}

	return c.redis.Set(
		ctx,
		c.buildRedisCooldownKey(channelId, eventName),
		"",
		time.Duration(seconds)*time.Second,
	).Err()
}
