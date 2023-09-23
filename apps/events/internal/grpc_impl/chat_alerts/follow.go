package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) Follow(ctx context.Context, req *events.FollowMessage) {
	if !c.settings.Followers.Enabled {
		return
	}

	if len(c.settings.Followers.Messages) == 0 {
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
}
