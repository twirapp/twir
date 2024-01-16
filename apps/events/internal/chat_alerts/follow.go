package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/bots"
	"github.com/satont/twir/libs/grpc/events"
)

const followCooldownKey = "follow"

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
