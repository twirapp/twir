package chat_alerts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *ChatAlerts) streamOnline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, platform.PlatformTwitch, text)
}

func (c *ChatAlerts) streamOffline(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
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

	return c.sendMessage(ctx, channelID, platform.PlatformTwitch, text)
}
