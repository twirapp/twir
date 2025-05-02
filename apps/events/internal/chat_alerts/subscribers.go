package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
)

type SubscribeMessage struct {
	UserName  string `json:"user_name"`
	Months    int    `json:"months"`
	ChannelId string `json:"channel_id"`
}

func (c *ChatAlerts) subscribe(
	ctx context.Context,
	settings model.ChatAlertsSettings,
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

	return c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      req.ChannelId,
			Message:        sample,
			SkipRateLimits: true,
		},
	)
}
