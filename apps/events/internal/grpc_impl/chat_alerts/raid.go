package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) Raid(ctx context.Context, msg *events.RaidedMessage) {
	if !c.settings.Raids.Enabled {
		return
	}

	sample := c.takeCountedSample(int(msg.Viewers), c.settings.Raids.Messages)
	sample = strings.ReplaceAll(sample, "{count}", fmt.Sprint(msg.Viewers))
	sample = strings.ReplaceAll(sample, "{user}", msg.UserName)

	if sample == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      msg.BaseInfo.ChannelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)
}
