package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
)

func (c *Handler) HandleChannelModeratorAdd(
	ctx context.Context,
	event eventsub.ChannelModeratorAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel moderator add",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserId, event.UserId, true)

	c.twirBus.Events.ModeratorAdded.Publish(
		ctx,
		events.ModeratorAddedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserId, event.UserId, true); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
}

func (c *Handler) HandleChannelModeratorRemove(
	ctx context.Context,
	event eventsub.ChannelModeratorRemoveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel moderator remove",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserId, event.UserId, false)

	c.twirBus.Events.ModeratorRemoved.Publish(
		ctx,
		events.ModeratorRemovedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserId, event.UserId, false); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
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
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if userId != channel.BotID {
		return
	}

	channel.IsBotMod = newStatus
	err = c.gorm.Save(&channel).Error
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	} else {
		if err = c.channelsCache.Invalidate(ctx, channelId); err != nil {
			c.logger.Error(err.Error(), logger.Error(err))
		}
	}
}
