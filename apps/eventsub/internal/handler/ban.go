package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/samber/lo"
	deprecatedmodel "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleBan(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelBan,
) {
	c.logger.Info(
		"channel ban",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
		slog.String("moderatorName", event.ModeratorUserName),
		slog.String("moderatorId", event.ModeratorUserID),
	)

	go func() {
		channel := deprecatedmodel.Channels{}
		if err := c.gorm.
			Where(`"id" = ?`, event.BroadcasterUserID).
			First(&channel).
			Error; err != nil {
			c.logger.Error("channel not found", slog.Any("err", err))
			return
		}

		if channel.BotID == event.UserID {
			channel.IsEnabled = false
			if err := c.gorm.Save(&channel).Error; err != nil {
				c.logger.Error("failed to disable channel", slog.Any("err", err))
			}

			return
		}
	}()

	t, _ := time.Parse(time.RFC3339, event.EndsAt)
	banEndsIn := t.Sub(time.Now().UTC())
	endsAt := lo.If(event.IsPermanent, "permanent").Else(
		fmt.Sprintf(
			"%v",
			math.Round(banEndsIn.Minutes()),
		),
	)

	c.twirBus.Events.ChannelBan.Publish(
		events.ChannelBanMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:             event.UserName,
			UserLogin:            event.UserLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserName,
			ModeratorUserName:    event.ModeratorUserName,
			ModeratorUserLogin:   event.ModeratorUserLogin,
			Reason:               event.Reason,
			EndsAt:               endsAt,
			IsPermanent:          event.IsPermanent,
		},
	)

	if err := c.eventsListRepository.Create(
		context.TODO(),
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
			Type:      model.ChannelEventListItemTypeChannelBan,
			Data: &model.ChannelsEventsListItemData{
				BanReason:            event.Reason,
				BanEndsInMinutes:     endsAt,
				BannedUserLogin:      event.UserLogin,
				BannedUserName:       event.UserName,
				ModeratorDisplayName: event.ModeratorUserName,
				ModeratorName:        event.ModeratorUserLogin,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	go c.websocketsGrpc.DudesUserPunished(
		context.TODO(),
		&websockets.DudesUserPunishedRequest{
			ChannelId:       event.BroadcasterUserID,
			UserId:          event.UserID,
			UserDisplayName: event.UserName,
			UserName:        event.UserLogin,
		},
	)
}
