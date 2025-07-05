package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func (c *ChatAlerts) streamOnline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req twitch.StreamOnlineMessage,
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
	text = strings.ReplaceAll(text, "{category}", req.CategoryName)

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.ChannelID,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}

func (c *ChatAlerts) streamOffline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req twitch.StreamOfflineMessage,
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

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.ChannelID,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
