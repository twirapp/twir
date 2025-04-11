package handler

import (
	"context"
	"log/slog"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelUpdate(
	_ *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelUpdate,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.logger.Info(
		"channel update",
		slog.String("title", event.Title),
		slog.String("category", event.CategoryName),
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	c.eventsGrpc.TitleOrCategoryChanged(
		ctx,
		&events.TitleOrCategoryChangedMessage{
			BaseInfo:    &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			NewTitle:    event.Title,
			NewCategory: event.CategoryName,
		},
	)

	if err := c.channelsInfoHistoryRepo.Create(
		ctx,
		channelsinfohistory.CreateInput{
			ChannelID: event.BroadcasterUserID,
			Title:     event.Title,
			Category:  event.CategoryName,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	err := c.gorm.
		WithContext(ctx).
		Model(&model.ChannelsStreams{}).
		Where(`"userId" = ?`, event.BroadcasterUserID).
		Updates(
			map[string]any{
				"title":      event.Title,
				`"gameName"`: event.CategoryName,
				`"gameId"`:   event.CategoryID,
			},
		).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
