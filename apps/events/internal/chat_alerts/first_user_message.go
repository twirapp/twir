package chat_alerts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) firstUserMessage(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, text)
}
