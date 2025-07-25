package handler

import (
	"context"
	"log/slog"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelModeratorAdd(
	ctx context.Context,
	_ *esb.ResponseHeaders, event *esb.EventChannelModeratorAdd,
) {
	c.logger.Info(
		"channel moderator add",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserID, event.UserID, true)

	c.twirBus.Events.ModeratorAdded.Publish(
		ctx,
		events.ModeratorAddedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserID,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserID, event.UserID, true); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) handleChannelModeratorRemove(
	ctx context.Context,
	_ *esb.ResponseHeaders, event *esb.EventChannelModeratorRemove,
) {
	c.logger.Info(
		"channel moderator remove",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserID, event.UserID, false)

	c.twirBus.Events.ModeratorRemoved.Publish(
		ctx,
		events.ModeratorRemovedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserID,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserID, event.UserID, false); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) updateUserModStatus(
	ctx context.Context,
	channelId string,
	userId string,
	newStatus bool,
) error {
	userStats := model.UsersStats{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"userId" = ? and "channelId" = ?`, userId, channelId).
		First(&userStats).Error; err != nil {
		return err
	}

	userStats.IsMod = newStatus

	return c.gorm.Save(&userStats).Error
}

func (c *Handler) updateBotStatus(
	ctx context.Context,
	channelId string,
	userId string,
	newStatus bool,
) {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", channelId).First(&channel).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	if userId != channel.BotID {
		return
	}

	channel.IsBotMod = newStatus
	err = c.gorm.Save(&channel).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
