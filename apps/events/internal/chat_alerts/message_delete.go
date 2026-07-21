package chat_alerts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) messageDelete(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, text)
}
