package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
)

const subscribeCooldownKey = "subscribe"

func (c *ChatAlerts) Subscribe(ctx context.Context, months int, userName, channelId string) {
	if !c.settings.Subscribers.Enabled {
		return
	}

	if len(c.settings.Subscribers.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		channelId,
		subscribeCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := c.takeCountedSample(months, c.settings.Subscribers.Messages)
	sample = strings.ReplaceAll(sample, "{months}", fmt.Sprint(months))
	sample = strings.ReplaceAll(sample, "{user}", userName)

	if sample == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      channelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	if c.settings.Subscribers.Cooldown != 0 {
		c.SetCooldown(ctx, channelId, subscribeCooldownKey, c.settings.Subscribers.Cooldown)
	}
}
