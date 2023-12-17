package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatAlerts) raid(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.RaidedMessage,
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

	_, err := c.botsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	return err
}
