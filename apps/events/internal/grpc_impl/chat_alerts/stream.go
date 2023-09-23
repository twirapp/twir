package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) StreamOnline(ctx context.Context, req *events.StreamOnlineMessage) {
	if !c.settings.StreamOnline.Enabled {
		return
	}

	if len(c.settings.StreamOnline.Messages) == 0 {
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
}

func (c *ChatAlerts) StreamOffline(ctx context.Context, req *events.StreamOfflineMessage) {
	if !c.settings.StreamOnline.Enabled {
		return
	}

	if len(c.settings.StreamOnline.Messages) == 0 {
		return
	}

	sample := lo.Sample(c.settings.StreamOnline.Messages)

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
}
