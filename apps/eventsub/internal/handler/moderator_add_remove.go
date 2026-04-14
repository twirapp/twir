package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	platform "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
)

func (c *Handler) HandleChannelModeratorAdd(
	ctx context.Context,
	event eventsub.ChannelModeratorAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel moderator add",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserId, true)

	c.twirBus.Events.ModeratorAdded.Publish(
		ctx,
		events.ModeratorAddedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserId, event.UserId, true); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
}

func (c *Handler) HandleChannelModeratorRemove(
	ctx context.Context,
	event eventsub.ChannelModeratorRemoveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel moderator remove",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)
	c.updateBotStatus(ctx, event.BroadcasterUserId, false)

	c.twirBus.Events.ModeratorRemoved.Publish(
		ctx,
		events.ModeratorRemovedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	if err := c.updateUserModStatus(ctx, event.BroadcasterUserId, event.UserId, false); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
}

func (c *Handler) updateUserModStatus(
	ctx context.Context,
	channelId string,
	userId string,
	newStatus bool,
) error {
	userStats := model.UsersStats{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"userId" = ? and "channelId" = ?`, userId, channelId).
		First(&userStats).Error; err != nil {
		return err
	}

	userStats.IsMod = newStatus

	return c.gorm.Save(&userStats).Error
}

func (c *Handler) updateBotStatus(
	ctx context.Context,
	channelId string,
	newStatus bool,
) {
	account, err := c.userPlatformAccountsRepo.GetByPlatformUserID(
		ctx,
		platform.PlatformTwitch,
		channelId,
	)
	if err != nil {
		if errors.Is(err, user_platform_accounts.ErrNotFound) {
			c.logger.Error("cannot find platform account", logger.Error(err))
		} else {
			c.logger.Error("cannot resolve broadcaster platform account", logger.Error(err))
		}
		return
	}

	channel, err := c.channelsRepo.GetByUserIDAndPlatform(ctx, account.UserID, platform.PlatformTwitch)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return
		}

		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	channel, err = c.channelsRepo.Update(ctx, channel.ID, channelsrepository.UpdateInput{IsBotMod: &newStatus})
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if err = c.channelsCache.Invalidate(ctx, channel.ID); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}
