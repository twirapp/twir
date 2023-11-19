package handler

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelRaid(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelRaid,
) {
	zap.S().Infow(
		"channel raid",
		"channelId", event.ToBroadcasterUserID,
		"channelName", event.ToBroadcasterUserName,
		"userId", event.FromBroadcasterUserID,
		"userName", event.FromBroadcasterUserLogin,
		"viewers", event.Viewers,
	)

	if err := c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.ToBroadcasterUserID,
			UserID:    event.FromBroadcasterUserID,
			Type:      model.ChannelEventListItemTypeRaided,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				RaidedViewersCount:    strconv.Itoa(event.Viewers),
				RaidedFromUserName:    event.FromBroadcasterUserLogin,
				RaidedFromDisplayName: event.FromBroadcasterUserName,
			},
		},
	).Error; err != nil {
		zap.S().Error(
			"cannot create raid entity",
			slog.String("channelId", event.ToBroadcasterUserID),
			slog.String("userId", event.FromBroadcasterUserID),
			zap.Error(err),
		)
	}

	if _, err := c.services.Grpc.Events.Raided(
		context.Background(),
		&events.RaidedMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.ToBroadcasterUserID},
			UserName:        event.FromBroadcasterUserLogin,
			UserDisplayName: event.FromBroadcasterUserName,
			Viewers:         int64(event.Viewers),
			UserId:          event.FromBroadcasterUserID,
		},
	); err != nil {
		zap.S().Error(
			"cannot fire raid event",
			slog.String("channelId", event.ToBroadcasterUserID),
			slog.String("userId", event.FromBroadcasterUserID),
			zap.Error(err),
		)
	}
}
