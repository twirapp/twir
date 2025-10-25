package handler

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	model "github.com/twirapp/twir/libs/gomodels"
)

func (c *Handler) HandleUserAuthorizationRevoke(
	ctx context.Context,
	event eventsub.UserAuthorizationRevokeEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"HandleUserAuthorizationRevoke",
		slog.String("user_id", event.UserId),
		slog.String("user_login", event.UserLogin),
	)

	if err := c.gorm.WithContext(ctx).Model(&model.Channels{}).
		Where("id = ?", event.UserId).
		Update(`"isBotMod"`, false).
		Update(`"isEnabled"`, false).Error; err != nil {
		c.logger.Error("failed to update channel", slog.Any("err", err))
	}

	user := &model.Users{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?", event.UserId).
		First(user).Error; err != nil {
		c.logger.Error("failed to get user", slog.Any("err", err))
	}

	if user.TokenID.Valid {
		if err := c.gorm.
			WithContext(ctx).
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
		if err := c.gorm.WithContext(ctx).Save(&user).Error; err != nil {
			c.logger.Error("failed to update user", slog.Any("err", err))
		}
	}
}
