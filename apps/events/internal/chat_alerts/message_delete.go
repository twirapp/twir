package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) messageDelete(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req events.ChannelMessageDeleteMessage,
) error {
	if !settings.MessageDelete.Enabled {
		return nil
	}

	if len(settings.MessageDelete.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.MessageDelete.Messages)

	text := strings.ReplaceAll(sample.Text, "{userName}", req.UserName)

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
