package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const redemptionsCooldownKey = "redemptions"

func (c *ChatAlerts) Redemption(ctx context.Context, req *events.RedemptionCreatedMessage) {
	if !c.settings.Redemptions.Enabled {
		return
	}

	if len(c.settings.Redemptions.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		"redemption",
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.Redemptions.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{user}", req.UserName)
	text = strings.ReplaceAll(text, "{reward}", req.RewardName)

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

	if c.settings.Redemptions.Cooldown != 0 {
		c.SetCooldown(
			ctx,
			req.BaseInfo.ChannelId,
			redemptionsCooldownKey,
			c.settings.Redemptions.Cooldown,
		)
	}
}
