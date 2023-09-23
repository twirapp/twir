package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) FirstUserMessage(ctx context.Context, req *events.FirstUserMessageMessage) {
	if !c.settings.FirstUserMessage.Enabled {
		return
	}

	if len(c.settings.FirstUserMessage.Messages) == 0 {
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
}
