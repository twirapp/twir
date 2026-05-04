package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/entities/platform"
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
		ctx,
		bots.SendMessageRequest{
			ChannelName:       lo.If(req.ChannelName != "", &req.ChannelName).Else(nil),
			ChannelId:         req.ChannelId,
			PlatformChannelID: req.ChannelId,
			Platform:          req.Platform.String(),
			Message:           sample,
			SkipRateLimits:    true,
		},
	)
}
