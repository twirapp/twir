package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func (c *Handler) resolveUserAndChannel(
	ctx context.Context,
	userPlatformID, broadcasterPlatformID string,
) (userID string, channelID string, err error) {
	chatUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, userPlatformID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("cannot resolve user: %w", err)
	}

	broadcasterUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, broadcasterPlatformID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("cannot resolve broadcaster user: %w", err)
	}

	channel, err := c.channelsRepo.GetByTwitchUserID(ctx, broadcasterUser.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("cannot get channel: %w", err)
	}

	return chatUser.ID.String(), channel.ID.String(), nil
}

func (c *Handler) resolveChannelIDByTwitchBroadcasterID(
	ctx context.Context,
	broadcasterPlatformID string,
) (string, error) {
	broadcasterUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, broadcasterPlatformID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return "", nil
		}

		return "", fmt.Errorf("cannot resolve broadcaster user: %w", err)
	}

	channel, err := c.channelsRepo.GetByTwitchUserID(ctx, broadcasterUser.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return "", nil
		}

		return "", fmt.Errorf("cannot get channel: %w", err)
	}

	return channel.ID.String(), nil
}

func (c *Handler) HandleChannelVipAdd(
	ctx context.Context,
	event eventsub.ChannelVipAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel vip add",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	c.twirBus.Events.VipAdded.Publish(
		ctx,
		events.VipAddedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	userID, channelID, err := c.resolveUserAndChannel(ctx, event.UserId, event.BroadcasterUserId)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
	if userID == "" || channelID == "" {
		return
	}

	if err := c.gorm.WithContext(ctx).Model(&model.UsersStats{}).
		Where(`"userId" = ?::uuid and "channelId" = ?::uuid`, userID, channelID).
		Update(`"is_vip"`, true).Error; err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}

func (c *Handler) HandleChannelVipRemove(
	ctx context.Context,
	event eventsub.ChannelVipRemoveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel vip remove",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	userID, channelID, err := c.resolveUserAndChannel(ctx, event.UserId, event.BroadcasterUserId)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
	if userID == "" || channelID == "" {
		return
	}

	if err := c.gorm.WithContext(ctx).Model(&model.UsersStats{}).
		Where(`"userId" = ?::uuid and "channelId" = ?::uuid`, userID, channelID).
		Update(`"is_vip"`, false).Error; err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.twirBus.Events.VipRemoved.Publish(
		ctx,
		events.VipRemovedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:   event.UserId,
			UserName: event.UserLogin,
		},
	)

	scheduledVip, err := c.scheduledVipsRepo.GetByUserAndChannelID(
		ctx,
		userID,
		channelID,
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if !scheduledVip.IsNil() {
		if err := c.scheduledVipsRepo.Delete(ctx, scheduledVip.ID); err != nil {
			c.logger.Error(err.Error(), logger.Error(err))
		}
	}
}
