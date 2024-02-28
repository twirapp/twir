package handler

import (
	"context"
	"log/slog"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *Handler) handleChannelUpdate(
	_ *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelUpdate,
) {
	c.logger.Info(
		"channel update",
		slog.String("title", event.Title),
		slog.String("category", event.CategoryName),
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	c.eventsGrpc.TitleOrCategoryChanged(
		context.Background(),
		&events.TitleOrCategoryChangedMessage{
			BaseInfo:    &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			NewTitle:    event.Title,
			NewCategory: event.CategoryName,
		},
	)

	err := c.gorm.Create(
		&model.ChannelInfoHistory{
			ID:        uuid.New().String(),
			Category:  event.CategoryName,
			Title:     event.Title,
			CreatedAt: time.Now().UTC(),
			ChannelID: event.BroadcasterUserID,
		},
	).Error

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	err = c.gorm.
		Model(&model.ChannelsStreams{}).
		Where(`"userId" = ?`, event.BroadcasterUserID).
		Updates(
			map[string]any{
				"title":    event.Title,
				"gameName": event.CategoryName,
			},
		).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
