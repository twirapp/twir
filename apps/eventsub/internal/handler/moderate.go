package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
)

func (c *Handler) HandleChannelModerateV2(
	ctx context.Context,
	event eventsub.ChannelModerateEventV2,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel moderate action",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("moderatorName", event.ModeratorUserName),
		slog.String("moderatorId", event.ModeratorUserID),
		slog.String("action", string(event.Action)),
	)

	switch event.Action {
	case eventsub.ChannelModerateEventV2ActionBan, eventsub.ChannelModerateEventV2ActionTimeout:
		c.handleModerateActionBan(ctx, event)
	case eventsub.ChannelModerateEventV2ActionUnban, eventsub.ChannelModerateEventV2ActionUntimeout:
		c.handleModerateActionUnBan(ctx, event)
	}
}
