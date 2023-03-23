package handler

import (
	"github.com/dnsge/twitch-eventsub-bindings"
	model "github.com/satont/tsuwari/libs/gomodels"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelModeratorAdd(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelModeratorAdd) {
	defer zap.S().Infow("channel moderator add",
		"channelId", event.BroadcasterUserID,
		"userId", event.UserID,
		"userName", event.UserLogin,
	)
	c.updateBotStatus(event.BroadcasterUserID, event.UserID, true)
}

func (c *Handler) handleChannelModeratorRemove(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelModeratorRemove) {
	defer zap.S().Infow("channel moderator remove",
		"channelId", event.BroadcasterUserID,
		"userId", event.UserID,
		"userName", event.UserLogin,
	)
	c.updateBotStatus(event.BroadcasterUserID, event.UserID, false)
}

func (c *Handler) updateBotStatus(channelId string, userId string, newStatus bool) {
	channel := model.Channels{}
	err := c.services.Gorm.Where("id = ?", channelId).First(&channel).Error
	if err != nil {
		zap.S().Errorw("failed to get channel", "error", err)
		return
	}

	if userId != channel.BotID {
		return
	}

	channel.IsBotMod = newStatus
	err = c.services.Gorm.Save(&channel).Error
	if err != nil {
		zap.S().Errorw("failed to update channel", "error", err)
	}
}
