package chat_alerts

import (
	"context"
	"strconv"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) Donation(ctx context.Context, msg *events.DonateMessage) {
	if !c.settings.Donations.Enabled {
		return
	}

	amount, err := strconv.Atoi(msg.Amount)
	if err != nil {
		return
	}

	sample := c.takeCountedSample(amount, c.settings.Donations.Messages)
	sample = strings.ReplaceAll(sample, "{count}", msg.Amount)
	sample = strings.ReplaceAll(sample, "{user}", msg.UserName)
	sample = strings.ReplaceAll(sample, "{message}", msg.Message)
	sample = strings.ReplaceAll(sample, "{currency}", msg.Currency)

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
