package chat_alerts

import (
	"context"
	"strconv"
	"strings"

	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

const donationsCooldownKey = "donations"

func (c *ChatAlerts) Donation(ctx context.Context, req *events.DonateMessage) {
	if !c.settings.Donations.Enabled {
		return
	}

	amount, err := strconv.Atoi(req.Amount)
	if err != nil {
		return
	}

	if cooldowned, err := c.IsOnCooldown(
		ctx,
		req.BaseInfo.ChannelId,
		donationsCooldownKey,
	); err != nil || cooldowned {
		return
	}

	sample := c.takeCountedSample(amount, c.settings.Donations.Messages)
	sample = strings.ReplaceAll(sample, "{count}", req.Amount)
	sample = strings.ReplaceAll(sample, "{user}", req.UserName)
	sample = strings.ReplaceAll(sample, "{message}", req.Message)
	sample = strings.ReplaceAll(sample, "{currency}", req.Currency)

	if sample == "" {
		return
	}

	c.services.BotsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	if c.settings.Donations.Cooldown != 0 {
		c.SetCooldown(ctx, req.BaseInfo.ChannelId, donationsCooldownKey, c.settings.Donations.Cooldown)
	}
}
