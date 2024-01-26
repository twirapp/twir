package handler

import (
	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
)

func (c *Handler) handleChannelChatMessage(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChatMessage,
) {
	
}
