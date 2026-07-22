package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/apps/eventsub/internal/channelbinding"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
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
				ChannelPlatformID: event.BroadcasterUserId,
				ChannelName:       event.BroadcasterUserLogin,
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
				ChannelPlatformID: event.BroadcasterUserId,
				ChannelName:       event.BroadcasterUserLogin,
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
	chatUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, userId)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("cannot resolve moderator user: %w", err)
	}

	broadcasterUser, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, channelId)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("cannot resolve broadcaster user: %w", err)
	}

	channel, err := c.channelService.GetChannelByBindingUserID(
		ctx,
		platform.PlatformTwitch,
		broadcasterUser.ID,
	)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("cannot get channel: %w", err)
	}

	userStats := model.UsersStats{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`user_id = ? and channel_id = ?`, chatUser.ID.String(), channel.ID.String()).
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
	user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, channelId)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			c.logger.Error("cannot find user by platform ID", logger.Error(err))
		} else {
			c.logger.Error("cannot resolve broadcaster user", logger.Error(err))
		}
		return
	}

	channel, err := c.channelService.GetChannelByBindingUserID(ctx, platform.PlatformTwitch, user.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return
		}

		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	twitchBinding, hasTwitchBinding := channelbinding.Find(channel, platform.PlatformTwitch)
	if !hasTwitchBinding {
		return
	}

	botConfigPatch := json.RawMessage(`{"is_bot_mod":false}`)
	if newStatus {
		botConfigPatch = json.RawMessage(`{"is_bot_mod":true}`)
	}
	if _, err = c.channelPlatformsRepo.Patch(ctx, twitchBinding.ID, channelplatforms.PatchInput{
		BotConfigPatch: botConfigPatch,
	}); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if err = c.channelsCache.Invalidate(ctx, channel.ID.String()); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}
