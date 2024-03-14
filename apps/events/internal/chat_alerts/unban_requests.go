package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *ChatAlerts) unbanRequestCreate(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.ChannelUnbanRequestCreateMessage,
) error {
	if !settings.UnbanRequestCreate.Enabled {
		return nil
	}

	if len(settings.UnbanRequestCreate.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.UnbanRequestCreate.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{userName}", req.UserName)
	text = strings.ReplaceAll(text, "{message}", req.Text)

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

func (c *ChatAlerts) unbanRequestResolved(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.ChannelUnbanRequestResolveMessage,
) error {
	if !settings.UnbanRequestResolve.Enabled {
		return nil
	}

	if len(settings.UnbanRequestResolve.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.UnbanRequestResolve.Messages)

	status := "approved"
	if req.Declined {
		status = "rejected"
	}

	text := sample.Text
	text = strings.ReplaceAll(text, "{userName}", req.UserName)
	text = strings.ReplaceAll(text, "{moderatorName}", req.ModeratorUserLogin)
	text = strings.ReplaceAll(text, "{message}", req.Reason)
	text = strings.ReplaceAll(text, "{status}", status)

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
