package chat_alerts

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const chatClearedCooldownKey = "chat_cleared"

func (c *ChatAlerts) ChatCleared(ctx context.Context, req *events.ChatClearMessage) {
	if !c.settings.ChatCleared.Enabled {
		return
	}

	if len(c.settings.ChatCleared.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		chatClearedCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.ChatCleared.Messages)

	if sample.Text == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample.Text,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	if c.settings.ChatCleared.Cooldown != 0 {
		c.SetCooldown(
			ctx,
			req.BaseInfo.ChannelId,
			chatClearedCooldownKey,
			c.settings.ChatCleared.Cooldown,
		)
	}
}
