package chat_alerts

import (
	"context"
	"strconv"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/grpc/events"
)

func (c *ChatAlerts) donation(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.DonateMessage,
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

	_, err = c.botsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)
	return err
}
