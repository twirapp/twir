package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const firstUserMessageCooldownKey = "first_user_message"

func (c *ChatAlerts) FirstUserMessage(ctx context.Context, req *events.FirstUserMessageMessage) {
	if !c.settings.FirstUserMessage.Enabled {
		return
	}

	if len(c.settings.FirstUserMessage.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		firstUserMessageCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.FirstUserMessage.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{user}", req.UserName)

	if text == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        text,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	if c.settings.FirstUserMessage.Cooldown != 0 {
		c.SetCooldown(
			ctx,
			req.BaseInfo.ChannelId,
			firstUserMessageCooldownKey,
			c.settings.FirstUserMessage.Cooldown,
		)
	}
}
