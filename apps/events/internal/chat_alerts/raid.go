package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *ChatAlerts) raid(
	ctx context.Context,
	settings model.ChatAlertsSettings,
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

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        sample,
			SkipRateLimits: true,
		},
	)
}
