package chat_alerts

import (
	"context"
	"strconv"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) donation(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req events.DonateMessage,
) error {
	if !settings.Donations.Enabled {
		return nil
	}

	amount, err := strconv.Atoi(req.Amount)
	if err != nil {
		return err
	}

	sample := c.takeCountedSample(amount, settings.Donations.Messages)
	sample = strings.ReplaceAll(sample, "{count}", req.Amount)
	sample = strings.ReplaceAll(sample, "{user}", req.UserName)
	sample = strings.ReplaceAll(sample, "{message}", req.Message)
	sample = strings.ReplaceAll(sample, "{currency}", req.Currency)

	if sample == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        sample,
			SkipRateLimits: true,
		},
	)
}
