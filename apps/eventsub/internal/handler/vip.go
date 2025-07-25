package handler

import (
	"context"
	"log/slog"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	scheduledmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelVipAdd(
	ctx context.Context,
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.ChannelVipAdd,
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelVipRemove(
	ctx context.Context,
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.ChannelVipRemove,
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
		c.logger.Error(err.Error(), slog.Any("err", err))
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
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	if scheduledVip != scheduledmodel.Nil {
		if err := c.scheduledVipsRepo.Delete(ctx, scheduledVip.ID); err != nil {
			c.logger.Error(err.Error(), slog.Any("err", err))
		}
	}
}
