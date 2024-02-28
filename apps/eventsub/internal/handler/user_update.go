package handler

import (
	"log/slog"

	"github.com/dnsge/twitch-eventsub-bindings"
)

func (c *Handler) handleUserUpdate(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventUserUpdate,
) {
	c.logger.Info(
		"user updated",
		slog.String("userId", event.UserID),
		slog.String("userLogin", event.UserLogin),
	)
}
