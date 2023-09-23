package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
)

func (c *ChatAlerts) Subscribe(ctx context.Context, months int, userName, channelId string) {
	if !c.settings.Subscribers.Enabled {
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
}
