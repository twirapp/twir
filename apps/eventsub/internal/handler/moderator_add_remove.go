package handler

import (
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelModeratorAdd(
	_ *esb.ResponseHeaders, event *esb.EventChannelModeratorAdd,
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
	_ *esb.ResponseHeaders, event *esb.EventChannelModeratorRemove,
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
