package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
)

func (c *Handler) HandleUserUpdate(
	ctx context.Context,
	event eventsub.UserUpdateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"user updated",
		slog.String("userId", event.UserId),
		slog.String("userLogin", event.UserLogin),
	)
}
