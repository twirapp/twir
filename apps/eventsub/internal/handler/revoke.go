package handler

import (
	"context"
	"database/sql"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleUserAuthorizationRevoke(
	ctx context.Context,
	h *esb.ResponseHeaders,
	event *esb.EventUserAuthorizationRevoke,
) {
	c.logger.Info(
		"handleUserAuthorizationRevoke",
		slog.String("user_id", event.UserID),
		slog.String("user_login", event.UserLogin),
	)

	if err := c.gorm.WithContext(ctx).Model(&model.Channels{}).
		Where("id = ?", event.UserID).
		Update(`"isBotMod"`, false).
		Update(`"isEnabled"`, false).Error; err != nil {
		c.logger.Error("failed to update channel", slog.Any("err", err))
	}

	user := &model.Users{}
	if err := c.gorm.
		WithContext(ctx).
		Where("id = ?", event.UserID).
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

func (c *Handler) handleSubRevocate(
	ctx context.Context,
	_ *esb.ResponseHeaders,
	revocation *esb.RevocationNotification,
) {
	c.logger.Info("handleSubRevocate", slog.Any("revocation", revocation))

	if revocation.Subscription.Status == "notification_failures_exceeded" {
		c.manager.SubscribeWithLimits(
			ctx,
			&eventsub_framework.SubRequest{
				Type:      revocation.Subscription.Type,
				Condition: revocation.Subscription.Condition,
				Callback:  revocation.Subscription.Transport.Callback,
				Secret:    c.config.TwitchClientSecret,
				Version:   revocation.Subscription.Version,
			},
		)
	}
}
