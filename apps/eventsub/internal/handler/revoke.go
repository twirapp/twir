package handler

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
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

	channel, channelErr := c.channelsRepo.GetByTwitchUserID(ctx, user.ID)
	if channelErr != nil {
		if !errors.Is(channelErr, channelsrepository.ErrNotFound) {
			c.logger.Error("failed to get channel", logger.Error(channelErr))
		}
	} else {
		isBotMod := false
		twitchEnabled := false
		overallEnabled := channel.KickBotJoined()

		if _, updateErr := c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{
			IsBotMod:         &isBotMod,
			IsEnabled:        &overallEnabled,
			TwitchBotEnabled: &twitchEnabled,
		}); updateErr != nil {
			c.logger.Error("failed to update channel", logger.Error(updateErr))
		}
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
