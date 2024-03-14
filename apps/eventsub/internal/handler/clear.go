package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelChatClear(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatClear,
) {
	c.logger.Info(
		"channel chat clear",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	c.gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			Type:      model.ChannelEventListItemTypeChatClear,
			CreatedAt: time.Now().UTC(),
			Data:      &model.ChannelsEventsListItemData{},
		},
	)

	c.eventsGrpc.ChatClear(
		context.Background(),
		&events.ChatClearMessage{
			BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		},
	)
}
