package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const followCooldownKey = "follow"

func (c *ChatAlerts) Follow(ctx context.Context, req *events.FollowMessage) {
	if !c.settings.Followers.Enabled {
		return
	}

	if len(c.settings.Followers.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		followCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.Followers.Messages)

	text := strings.ReplaceAll(sample.Text, "{user}", req.UserName)

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

	if c.settings.Followers.Cooldown != 0 {
		c.SetCooldown(ctx, req.BaseInfo.ChannelId, followCooldownKey, c.settings.Followers.Cooldown)
	}
}
