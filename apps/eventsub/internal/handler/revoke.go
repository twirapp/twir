package handler

import (
	"database/sql"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleUserAuthorizationRevoke(
	h *esb.ResponseHeaders,
	event *esb.EventUserAuthorizationRevoke,
) {
	c.logger.Info(
		"handleUserAuthorizationRevoke",
		slog.String("user_id", event.UserID),
		slog.String("user_login", event.UserLogin),
	)

	if err := c.gorm.Model(&model.Channels{}).
		Where("id = ?", event.UserID).
		Update(`"isBotMod"`, false).
		Update(`"isEnabled"`, false).Error; err != nil {
		c.logger.Error("failed to update channel", slog.Any("err", err))
	}

	user := &model.Users{}
	if err := c.gorm.
		Where("id = ?", event.UserID).
		First(user).Error; err != nil {
		c.logger.Error("failed to get user", slog.Any("err", err))
	}

	if user.TokenID.Valid {
		if err := c.gorm.
			Delete(
				&model.Tokens{},
				"id = ?",
				user.TokenID.String,
			).Error; err != nil {
			c.logger.Error(
				"failed to delete token",
				slog.Any("err", err),
			)
		}

		user.TokenID = sql.NullString{}
		if err := c.gorm.Save(&user).Error; err != nil {
			c.logger.Error("failed to update user", slog.Any("err", err))
		}
	}
}

func (c *Handler) handleSubRevocate(
	h *esb.ResponseHeaders,
	revocation *esb.RevocationNotification,
) {
	topic := &model.EventsubTopic{}
	if err := c.gorm.
		Where("topic = ?", revocation.Subscription.Type).
		First(topic).Error; err != nil {
		c.logger.Error("failed to get topic", slog.Any("err", err))
		return
	}

	subscription := &model.EventsubSubscription{}
	if err := c.gorm.
		Where("topic_id = ?", topic.ID).
		First(subscription).Error; err != nil {
		c.logger.Error("failed to get subscription", slog.Any("err", err))
		return
	}

	subscription.Status = revocation.Subscription.Status

	if err := c.gorm.Save(&subscription).Error; err != nil {
		c.logger.Error("failed to delete subscription", slog.Any("err", err))
	}
}
