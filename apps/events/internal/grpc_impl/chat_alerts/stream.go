package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const streamOnlineCooldownKey = "stream_online"
const streamOfflineCooldownKey = "stream_offline"

func (c *ChatAlerts) StreamOnline(ctx context.Context, req *events.StreamOnlineMessage) {
	if !c.settings.StreamOnline.Enabled {
		return
	}

	if len(c.settings.StreamOnline.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		streamOnlineCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.StreamOnline.Messages)
	text := sample.Text
	text = strings.ReplaceAll(text, "{title}", req.Title)
	text = strings.ReplaceAll(text, "{category}", req.Category)

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

	if c.settings.StreamOnline.Cooldown != 0 {
		c.SetCooldown(
			ctx,
			req.BaseInfo.ChannelId,
			streamOnlineCooldownKey,
			c.settings.StreamOnline.Cooldown,
		)
	}
}

func (c *ChatAlerts) StreamOffline(ctx context.Context, req *events.StreamOfflineMessage) {
	if !c.settings.StreamOffline.Enabled {
		return
	}

	if len(c.settings.StreamOffline.Messages) == 0 {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		streamOfflineCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := lo.Sample(c.settings.StreamOffline.Messages)

	text := sample.Text
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

	if c.settings.StreamOffline.Cooldown != 0 {
		c.SetCooldown(
			ctx,
			req.BaseInfo.ChannelId,
			streamOfflineCooldownKey,
			c.settings.StreamOffline.Cooldown,
		)
	}
}
