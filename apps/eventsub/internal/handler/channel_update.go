package handler

import (
	"context"
	"encoding/json"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/zap"
	"time"
)

func (c *Handler) handleChannelUpdate(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelUpdate,
) {
	defer zap.S().Infow(
		"channel update",
		"title", event.Title,
		"category", event.CategoryName,
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
	)

	pbMessage := pubsub.StreamUpdateMessage{
		ChannelID: event.BroadcasterUserID,
		Title:     event.Title,
		Category:  event.CategoryName,
	}
	bytes, err := json.Marshal(pbMessage)
	if err != nil {
		zap.S().Errorw("failed to marshal message", "error", err)
		return
	}
	c.services.PubSub.Publish("stream.update", bytes)

	c.services.Grpc.Events.TitleOrCategoryChanged(
		context.Background(), &events.TitleOrCategoryChangedMessage{
			BaseInfo:    &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			NewTitle:    event.Title,
			NewCategory: event.CategoryName,
		},
	)

	err = c.services.Gorm.Create(
		&model.ChannelInfoHistory{
			ID:        uuid.New().String(),
			Category:  event.CategoryName,
			Title:     event.Title,
			CreatedAt: time.Now().UTC(),
			ChannelID: event.BroadcasterUserID,
		},
	).Error

	if err != nil {
		zap.S().Errorw("failed to save channel info history", "error", err)
	}

	err = c.services.Gorm.
		Model(&model.ChannelsStreams{}).
		Where(`"userId" = ?`, event.BroadcasterUserID).
		Updates(
			map[string]any{
				"title":    event.Title,
				"gameName": event.CategoryName,
			},
		).Error
	if err != nil {
		zap.S().Errorw("failed to update channel stream", "error", err)
	}
}
