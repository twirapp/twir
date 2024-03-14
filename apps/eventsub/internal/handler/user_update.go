package handler

import (
	"log/slog"

	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleUserUpdate(
	_ *esb.ResponseHeaders,
	event *esb.EventUserUpdate,
) {
	c.logger.Info(
		"user updated",
		slog.String("userId", event.UserID),
		slog.String("userLogin", event.UserLogin),
	)
}
