package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

type SubscribeMessage struct {
	UserName    string            `json:"user_name"`
	Months      int               `json:"months"`
	ChannelId   string            `json:"channel_id"`
	ChannelName string            `json:"channel_name"`
	Platform    platform.Platform `json:"platform"`
}

func (c *ChatAlerts) subscribe(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	channelID uuid.UUID,
	req SubscribeMessage,
) error {
	if !settings.Subscribers.Enabled {
		return nil
	}

	if len(settings.Subscribers.Messages) == 0 {
		return nil
	}

	sample := c.takeCountedSample(req.Months, settings.Subscribers.Messages)
	sample = strings.ReplaceAll(sample, "{month}", fmt.Sprint(req.Months))
	sample = strings.ReplaceAll(sample, "{months}", fmt.Sprint(req.Months))
	sample = strings.ReplaceAll(sample, "{user}", req.UserName)

	if sample == "" {
		return nil
	}

	return c.sendMessage(ctx, channelID, req.Platform, sample)
}
