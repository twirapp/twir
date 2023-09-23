package chat_alerts

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) ChatCleared(ctx context.Context, req *events.ChatClearMessage) {
	if !c.settings.ChatCleared.Enabled {
		return
	}

	if len(c.settings.ChatCleared.Messages) == 0 {
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
}
