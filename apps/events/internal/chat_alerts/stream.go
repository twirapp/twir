package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/grpc/events"
)

func (c *ChatAlerts) streamOnline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.StreamOnlineMessage,
) error {
	if !settings.StreamOnline.Enabled {
		return nil
	}

	if len(settings.StreamOnline.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.StreamOnline.Messages)
	text := sample.Text
	text = strings.ReplaceAll(text, "{title}", req.Title)
	text = strings.ReplaceAll(text, "{category}", req.Category)

	if text == "" {
		return nil
	}

	_, err := c.botsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        text,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)
	return err
}

func (c *ChatAlerts) streamOffline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.StreamOfflineMessage,
) error {
	if !settings.StreamOffline.Enabled {
		return nil
	}

	if len(settings.StreamOffline.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.StreamOffline.Messages)

	text := sample.Text
	if text == "" {
		return nil
	}

	_, err := c.botsGrpc.SendMessage(
		ctx, &bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        text,
			IsAnnounce:     nil,
			SkipRateLimits: true,
		},
	)

	return err
}
