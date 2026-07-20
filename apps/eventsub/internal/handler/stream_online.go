package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twir/libs/twitch"
)

func (c *Handler) HandleStreamOnline(
	ctx context.Context,
	event eventsub.StreamOnlineEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(event.BroadcasterUserId),
	).Err(); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.logger.Info(
		"stream online",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.config, c.twirBus)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	streamsReq, err := twitchClient.GetStreams(
		&helix.StreamsParams{
			UserIDs: []string{event.BroadcasterUserId},
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}

	if streamsReq.ErrorMessage != "" {
		c.logger.Error(streamsReq.ErrorMessage)
		return
	}

	i := 0
	for {
		if i > 5 {
			break
		}

		if len(streamsReq.Data.Streams) == 0 {
			c.logger.Error(
				"stream online event received but GetStreams returned no streams",
				slog.String("channelId", event.BroadcasterUserId),
				slog.String("channelName", event.BroadcasterUserLogin),
			)
			i++
			time.Sleep(5 * time.Second)
			continue
		}

		stream := streamsReq.Data.Streams[0]

		c.gorm.WithContext(ctx).Where(
			`"userId" = ?`,
			event.BroadcasterUserId,
		).Delete(&model.ChannelsStreams{})

		tags := pq.StringArray{}
		for _, tag := range stream.Tags {
			tags = append(tags, tag)
		}
		tagIds := pq.StringArray{}
		for _, tagId := range stream.TagIDs {
			tagIds = append(tagIds, tagId)
		}

		err = c.gorm.WithContext(ctx).Create(
			&model.ChannelsStreams{
				ID:           event.Id,
				UserId:       event.BroadcasterUserId,
				UserLogin:    event.BroadcasterUserLogin,
				UserName:     event.BroadcasterUserName,
				GameId:       stream.GameID,
				GameName:     stream.GameName,
				CommunityIds: nil,
				Type:         stream.Type,
				Title:        stream.Title,
				ViewerCount:  stream.ViewerCount,
				StartedAt:    stream.StartedAt,
				Language:     stream.Language,
				ThumbnailUrl: stream.ThumbnailURL,
				TagIds:       &tagIds,
				Tags:         &tags,
				IsMature:     stream.IsMature,
			},
		).Error
		if err != nil {
			c.logger.Error(
				"cannot create stream record",
				slog.String("channelId", event.BroadcasterUserId),
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
				StartedAt:    stream.StartedAt,
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
