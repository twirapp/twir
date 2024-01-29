package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/pubsub"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleStreamOnline(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventStreamOnline,
) {
	c.logger.Info(
		"stream online",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	twitchClient, err := twitch.NewAppClient(c.config, c.tokensGrpc)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	streamsReq, err := twitchClient.GetStreams(
		&helix.StreamsParams{
			UserIDs: []string{event.BroadcasterUserID},
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	if streamsReq.ErrorMessage != "" {
		c.logger.Error(streamsReq.ErrorMessage)
		return
	}

	if len(streamsReq.Data.Streams) == 0 {
		return
	}

	stream := streamsReq.Data.Streams[0]

	err = c.gorm.Where(
		`"userId" = ?`,
		event.BroadcasterUserID,
	).Delete(&model.ChannelsStreams{}).Error
	if err == nil {
		tags := pq.StringArray{}
		for _, tag := range stream.Tags {
			tags = append(tags, tag)
		}
		tagIds := pq.StringArray{}
		for _, tagId := range stream.TagIDs {
			tagIds = append(tagIds, tagId)
		}

		err = c.gorm.Create(
			&model.ChannelsStreams{
				ID:             event.ID,
				UserId:         event.BroadcasterUserID,
				UserLogin:      event.BroadcasterUserLogin,
				UserName:       event.BroadcasterUserName,
				GameId:         stream.GameID,
				GameName:       stream.GameName,
				CommunityIds:   nil,
				Type:           stream.Type,
				Title:          stream.Title,
				ViewerCount:    stream.ViewerCount,
				StartedAt:      stream.StartedAt,
				Language:       stream.Language,
				ThumbnailUrl:   stream.ThumbnailURL,
				TagIds:         &tagIds,
				Tags:           &tags,
				IsMature:       stream.IsMature,
				ParsedMessages: 0,
			},
		).Error
		if err != nil {
			zap.S().Error(err)
		}
	}

	_, err = c.eventsGrpc.StreamOnline(
		context.Background(),
		&events.StreamOnlineMessage{
			BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			Title:    streamsReq.Data.Streams[0].Title,
			Category: streamsReq.Data.Streams[0].GameName,
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	bytes, err := json.Marshal(
		&pubsub.StreamOnlineMessage{
			ChannelID: event.BroadcasterUserID,
			StreamID:  event.ID,
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	c.pubSub.Client.Publish("stream.online", bytes)
}
