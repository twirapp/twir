package handler

import (
	"context"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelChatClear(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatClear,
) {
	zap.S().Infow(
		"channel chat clear",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
	)

	c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			Type:      model.ChannelEventListItemTypeChatClear,
			CreatedAt: time.Now().UTC(),
			Data:      &model.ChannelsEventsListItemData{},
		},
	)
	c.services.Grpc.Events.ChatClear(
		context.Background(),
		&events.ChatClearMessage{
			BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		},
	)
}
