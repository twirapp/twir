package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) firstUserMessage(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req events.FirstUserMessageMessage,
) error {
	if !settings.FirstUserMessage.Enabled {
		return nil
	}

	if len(settings.FirstUserMessage.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.FirstUserMessage.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{user}", req.UserName)

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
