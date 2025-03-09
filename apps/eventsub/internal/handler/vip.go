package handler

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelVipAdd(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.ChannelVipAdd,
) {
	c.logger.Info(
		"channel vip add",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	if err := c.gorm.Model(&model.UsersStats{}).
		Where(`"userId" = ? and "channelId" = ?`, event.UserId, event.BroadcasterUserId).
		Update(`"is_vip"`, true).Error; err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelVipRemove(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.ChannelVipRemove,
) {
	c.logger.Info(
		"channel vip remove",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	if err := c.gorm.Model(&model.UsersStats{}).
		Where(`"userId" = ? and "channelId" = ?`, event.UserId, event.BroadcasterUserId).
		Update(`"is_vip"`, false).Error; err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
