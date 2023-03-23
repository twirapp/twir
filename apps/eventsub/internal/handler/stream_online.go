package handler

import (
	"context"
	"encoding/json"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/lib/pq"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"github.com/satont/tsuwari/libs/pubsub"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
)

func (c *Handler) handleStreamOnline(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventStreamOnline) {
	defer zap.S().Infow("stream online",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
	)

	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		zap.S().Error(err)
		return
	}

	streamsReq, err := twitchClient.GetStreams(&helix.StreamsParams{
		UserIDs: []string{event.BroadcasterUserID},
	})

	if err != nil || streamsReq.ErrorMessage != "" {
		zap.S().Error(err, streamsReq.ErrorMessage)
		return
	}

	if len(streamsReq.Data.Streams) == 0 {
		return
	}

	stream := streamsReq.Data.Streams[0]

	err = c.services.Gorm.Where(`"userId" = ?`, event.BroadcasterUserID).Delete(&model.ChannelsStreams{}).Error
	if err == nil {
		tags := pq.StringArray{}
		for _, tag := range stream.Tags {
			tags = append(tags, tag)
		}
		tagIds := pq.StringArray{}
		for _, tagId := range stream.TagIDs {
			tagIds = append(tagIds, tagId)
		}

		err = c.services.Gorm.Create(&model.ChannelsStreams{
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
		}).Error
		if err != nil {
			zap.S().Error(err)
		}
	}

	c.services.Grpc.Events.StreamOnline(context.Background(), &events.StreamOnlineMessage{
		BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		Title:    streamsReq.Data.Streams[0].Title,
		Category: streamsReq.Data.Streams[0].GameName,
	})

	bytes, err := json.Marshal(&pubsub.StreamOnlineMessage{
		ChannelID: event.BroadcasterUserID,
		StreamID:  event.ID,
	})
	if err != nil {
		zap.S().Error(err)
		return
	}

	c.services.PubSub.Publish("stream.online", bytes)
}
