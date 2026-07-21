package chat_alerts

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) chatCleared(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, sample.Text)
}
