package handler

import (
	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"go.uber.org/zap"
)

func (c *handler) handleChannelUpdate(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelUpdate) {
	defer zap.S().Infow("channel update",
		"title", event.Title,
		"category", event.CategoryName,
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
	)
}
