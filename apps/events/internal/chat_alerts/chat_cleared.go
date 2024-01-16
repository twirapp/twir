package chat_alerts

import (
	"context"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/grpc/events"
)

func (c *ChatAlerts) chatCleared(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.ChatClearMessage,
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

	_, err := c.botsGrpc.SendMessage(
		ctx,
		&bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        sample.Text,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	return err
}
