package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	model "github.com/twirapp/twir/libs/gomodels"
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

	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(event.BroadcasterUserId),
	).Err(); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	dbStream := model.ChannelsStreams{}
	if err := c.gorm.WithContext(ctx).Where(
		`"userId" = ?`,
		event.BroadcasterUserId,
	).First(&dbStream).Error; err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	c.twirBus.Channel.StreamOffline.Publish(
		ctx,
		twitch.StreamOfflineMessage{
			ChannelID: event.BroadcasterUserId,
			StartedAt: dbStream.StartedAt,
		},
	)

	err := c.gorm.WithContext(ctx).Where(
		`"userId" = ?`,
		event.BroadcasterUserId,
	).Delete(&model.ChannelsStreams{}).Error
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
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

	stream, err := c.streamsrepository.GetByChannelID(ctx, event.BroadcasterUserId)
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
			ChannelID:  &event.BroadcasterUserId,
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
		event.BroadcasterUserId,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return fmt.Errorf("failed to create twitch client: %w", err)
	}

	var wg sync.WaitGroup

	for _, vip := range vips {
		wg.Add(1)

		go func() {
			defer wg.Done()
			resp, err := twitchClient.RemoveChannelVip(
				&helix.RemoveChannelVipParams{
					BroadcasterID: vip.ChannelID,
					UserID:        vip.UserID,
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
		}()

	}

	wg.Wait()

	return nil
}
