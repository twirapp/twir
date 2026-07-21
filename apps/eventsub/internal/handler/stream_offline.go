package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/eventsub"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/entities/platform"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	twitchlib "github.com/twirapp/twir/libs/twitch"
)

func (c *Handler) HandleStreamOffline(
	ctx context.Context,
	event eventsub.StreamOfflineEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"stream offline",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	channel, err := c.channelService.GetChannelByPlatformUserID(
		ctx,
		event.BroadcasterUserId,
		platform.PlatformTwitch,
	)
	if err != nil {
		c.logger.Error(
			"cannot resolve channel for stream offline",
			slog.String("channelId", event.BroadcasterUserId),
			logger.Error(err),
		)
		return
	}

	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(channel.ID.String()),
	).Err(); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	dbStream, err := c.streamsrepository.GetByChannelID(ctx, channel.ID, platform.PlatformTwitch)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if dbStream.IsNil() {
		c.logger.Error(
			"stream offline event received but no stream record found",
			slog.String("channelId", event.BroadcasterUserId),
		)
		return
	}

	c.twirBus.Channel.StreamOffline.Publish(
		ctx,
		twitch.StreamOfflineMessage{
			ChannelID: event.BroadcasterUserId,
			StartedAt: dbStream.StartedAt,
		},
	)

	if err := c.streamsrepository.DeleteByChannelID(ctx, channel.ID, platform.PlatformTwitch); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	if err := c.channelService.InvalidateOnlineCache(ctx, channel.ID); err != nil {
		c.logger.Error(
			"cannot invalidate online cache",
			slog.String("channelId", channel.ID.String()),
			logger.Error(err),
		)
	}

	go func() {
		time.Sleep(10 * time.Minute)
		if err := c.handleStreamOfflineScheduledVips(event); err != nil {
			c.logger.Error(err.Error(), logger.Error(err))
		}
	}()
}

func (c *Handler) handleStreamOfflineScheduledVips(
	event eventsub.StreamOfflineEvent,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, event.BroadcasterUserId)
	if err != nil {
		return fmt.Errorf("failed to get user by platform id: %w", err)
	}

	channel, err := c.channelService.GetChannelByConnectedUser(ctx, user.ID, platform.PlatformTwitch)
	if err != nil {
		return fmt.Errorf("failed to get channel by broadcaster user: %w", err)
	}

	channelID := channel.ID.String()

	stream, err := c.streamsrepository.GetByChannelID(ctx, channel.ID, platform.PlatformTwitch)
	if err != nil {
		return fmt.Errorf("failed to get stream by channel id: %w", err)
	}

	if !stream.IsNil() {
		// Stream is live again, do not remove VIPs
		return nil
	}

	vips, err := c.scheduledVipsRepo.GetMany(
		ctx,
		scheduledvipsrepository.GetManyInput{
			ChannelID:  &channelID,
			RemoveType: lo.ToPtr(scheduledvipsentity.RemoveTypeStreamEnd),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get scheduled vips: %w", err)
	}

	if len(vips) == 0 {
		return nil
	}

	twitchClient, err := twitchlib.NewUserClientWithContext(
		ctx,
		user.ID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return fmt.Errorf("failed to create twitch client: %w", err)
	}

	var wg sync.WaitGroup

	for _, vip := range vips {
		wg.Add(1)

		go func(vip scheduledvipsentity.ScheduledVip) {
			defer wg.Done()

			vipUserID, err := uuid.Parse(vip.UserID)
			if err != nil {
				c.logger.Error(
					"failed to parse vip user id",
					slog.String("userId", vip.UserID),
					logger.Error(err),
				)
				return
			}

			vipUser, err := c.usersRepo.GetByID(ctx, vipUserID)
			if err != nil {
				c.logger.Error(
					"failed to get vip user by id",
					slog.String("userId", vip.UserID),
					logger.Error(err),
				)
				return
			}
			if vipUser.IsNil() {
				c.logger.Error("vip user not found", slog.String("userId", vip.UserID))
				return
			}

			resp, err := twitchClient.RemoveChannelVip(
				&helix.RemoveChannelVipParams{
					BroadcasterID: event.BroadcasterUserId,
					UserID:        vipUser.PlatformID,
				},
			)
			if err != nil {
				c.logger.Error("failed to remove vip", logger.Error(err))
				return
			}

			if resp.ErrorMessage != "" {
				c.logger.Error(
					"failed to remove vip",
					logger.Error(errors.New(resp.ErrorMessage)),
				)
				return
			}

			c.logger.Info(
				"removed vip on stream end",
				slog.String("userId", vip.UserID),
				slog.String("channelId", vip.ChannelID),
			)
		}(vip)

	}

	wg.Wait()

	return nil
}
