package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) raid(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
	req events.RaidedMessage,
) error {
	if !settings.Raids.Enabled {
		return nil
	}

	if len(settings.Raids.Messages) == 0 {
		return nil
	}

	sample := c.takeCountedSample(int(req.Viewers), settings.Raids.Messages)
	sample = strings.ReplaceAll(sample, "{count}", fmt.Sprint(req.Viewers))
	sample = strings.ReplaceAll(sample, "{user}", req.UserName)

	if sample == "" {
		return nil
	}

	return c.sendMessage(ctx, channelID, req.BaseInfo.Platform, sample)
}
