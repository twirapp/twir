package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/kvizyx/twitchy/helix"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *Handler) HandleStreamOnline(
	ctx context.Context,
	event eventsub.StreamOnlineEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platform.PlatformTwitch,
		event.BroadcasterUserId,
	)
	if err != nil {
		c.logger.Error(
			"cannot resolve channel for stream online",
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

	c.logger.Info(
		"stream online",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, nil)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	streamsReq, err := twitchClient.Streams.GetStreams(
		ctx,
		helix.GetStreamsRequest{
			UserID: []string{event.BroadcasterUserId},
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	i := 0
	for {
		if i > 5 {
			break
		}

		if len(streamsReq.Data) == 0 {
			c.logger.Error(
				"stream online event received but GetStreams returned no streams",
				slog.String("channelId", event.BroadcasterUserId),
				slog.String("channelName", event.BroadcasterUserLogin),
			)
			i++
			time.Sleep(5 * time.Second)
			continue
		}

		stream := streamsReq.Data[0]

		err = c.streamsrepository.Save(
			ctx,
			streamsrepository.SaveInput{
				ID:           event.Id,
				ChannelID:    channel.ID,
				UserId:       event.BroadcasterUserId,
				UserLogin:    event.BroadcasterUserLogin,
				UserName:     event.BroadcasterUserName,
				GameId:       stream.GameID,
				GameName:     stream.GameName,
				CommunityIds: nil,
				Type:         string(stream.Type),
				Title:        stream.Title,
				ViewerCount:  stream.ViewerCount,
				StartedAt:    stream.StartedAt.Time,
				Language:     stream.Language,
				ThumbnailUrl: stream.ThumbnailURL,
				TagIds:       stream.TagIDs,
				Tags:         stream.Tags,
				IsMature:     stream.IsMature,
				Platform:     platform.PlatformTwitch,
			},
		)
		if err != nil {
			c.logger.Error(
				"cannot create stream record",
				slog.String("channelId", event.BroadcasterUserId),
				logger.Error(err),
			)
		}

		if err := c.channelService.InvalidateOnlineCache(ctx, channel.ID); err != nil {
			c.logger.Error(
				"cannot invalidate online cache",
				slog.String("channelId", channel.ID.String()),
				logger.Error(err),
			)
		}

		c.twirBus.Channel.StreamOnline.Publish(
			ctx,
			bustwitch.StreamOnlineMessage{
				ChannelID:    event.BroadcasterUserId,
				StreamID:     event.Id,
				CategoryName: stream.GameName,
				CategoryID:   stream.GameID,
				Title:        stream.Title,
				Viewers:      stream.ViewerCount,
				StartedAt:    stream.StartedAt.Time,
			},
		)

		break
	}

	// c.channelsInfoHistoryRepo.Create(
	// 	ctx,
	// 	channelsinfohistory.CreateInput{
	// 		ChannelID: stream.UserID,
	// 		Title:     stream.Title,
	// 		Category:  stream.GameName,
	// 	},
	// )
}
