package chat_alerts

import (
	"context"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) chatCleared(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req events.ChatClearMessage,
) error {
	if !settings.ChatCleared.Enabled {
		return nil
	}

	if len(settings.ChatCleared.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.ChatCleared.Messages)

	if sample.Text == "" {
		return nil
	}

	err := c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        sample.Text,
			SkipRateLimits: true,
		},
	)

	return err
}
