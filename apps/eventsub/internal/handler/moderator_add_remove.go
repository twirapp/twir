package handler

import (
	"log/slog"

	"github.com/dnsge/twitch-eventsub-bindings"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Handler) handleChannelModeratorAdd(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelModeratorAdd,
) {
	c.logger.Info(
		"channel moderator add",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(event.BroadcasterUserID, event.UserID, true)
}

func (c *Handler) handleChannelModeratorRemove(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelModeratorRemove,
) {
	c.logger.Info(
		"channel moderator remove",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(event.BroadcasterUserID, event.UserID, false)
}

func (c *Handler) updateBotStatus(channelId string, userId string, newStatus bool) {
	channel := model.Channels{}
	err := c.gorm.Where("id = ?", channelId).First(&channel).Error
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
