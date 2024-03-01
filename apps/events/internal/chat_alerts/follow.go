package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *ChatAlerts) follow(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.FollowMessage,
) error {
	if !settings.Followers.Enabled {
		return nil
	}

	if len(settings.Followers.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.Followers.Messages)

	text := strings.ReplaceAll(sample.Text, "{user}", req.UserName)

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
