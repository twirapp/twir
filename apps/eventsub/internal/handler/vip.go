package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	scheduledmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
)

func (c *Handler) HandleChannelVipAdd(
	ctx context.Context,
	event eventsub.ChannelVipAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel vip add",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	c.twirBus.Events.VipAdded.Publish(
		ctx,
		events.VipAddedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	if err := c.gorm.WithContext(ctx).Model(&model.UsersStats{}).
		Where(`"userId" = ? and "channelId" = ?`, event.UserId, event.BroadcasterUserId).
		Update(`"is_vip"`, true).Error; err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}

func (c *Handler) HandleChannelVipRemove(
	ctx context.Context,
	event eventsub.ChannelVipRemoveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel vip remove",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	if err := c.gorm.WithContext(ctx).Model(&model.UsersStats{}).
		Where(`"userId" = ? and "channelId" = ?`, event.UserId, event.BroadcasterUserId).
		Update(`"is_vip"`, false).Error; err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.twirBus.Events.VipRemoved.Publish(
		ctx,
		events.VipRemovedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	scheduledVip, err := c.scheduledVipsRepo.GetByUserAndChannelID(
		ctx,
		event.UserId,
		event.BroadcasterUserId,
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if scheduledVip != scheduledmodel.Nil {
		if err := c.scheduledVipsRepo.Delete(ctx, scheduledVip.ID); err != nil {
			c.logger.Error(err.Error(), logger.Error(err))
		}
	}
}
