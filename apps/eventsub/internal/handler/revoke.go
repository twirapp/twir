package handler

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
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

	user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, event.UserId)
	if err != nil {
		c.logger.Error("failed to get user", logger.Error(err))
		return
	}

	if err := c.gorm.WithContext(ctx).Model(&model.Channels{}).
		Where("twitch_user_id = ?", user.ID).
		Update(`"isBotMod"`, false).
		Update(`"isEnabled"`, false).Error; err != nil {
		c.logger.Error("failed to update channel", logger.Error(err))
	}

	legacyUser := &model.Users{ID: user.ID.String(), TokenID: user.TokenID.NullString}
	if legacyUser.TokenID.Valid {
		if err := c.gorm.
			WithContext(ctx).
			Delete(
				&model.Tokens{},
				"id = ?",
				legacyUser.TokenID.String,
			).Error; err != nil {
			c.logger.Error(
				"failed to delete token",
				logger.Error(err),
			)
		}

		legacyUser.TokenID = sql.NullString{}
		if err := c.gorm.WithContext(ctx).Save(&legacyUser).Error; err != nil {
			c.logger.Error("failed to update user", logger.Error(err))
		}
	}
}
