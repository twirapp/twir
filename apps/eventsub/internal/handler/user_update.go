package handler

import (
	"encoding/json"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/zap"
)

func (c *Handler) handleUserUpdate(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventUserUpdate) {
	defer zap.S().Infow("user update", "userId", event.UserID, "userLogin", event.UserLogin)

	bytes, err := json.Marshal(
		&pubsub.UserUpdateMessage{
			UserID:      event.UserID,
			UserLogin:   event.UserLogin,
			UserName:    event.UserName,
			Description: event.Description,
		},
	)
	if err != nil {
		zap.S().Error(err)
	}

	c.services.PubSub.Publish("user.update", bytes)
}
