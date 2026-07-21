package chat_alerts

import (
	"context"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) donation(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, sample)
}
