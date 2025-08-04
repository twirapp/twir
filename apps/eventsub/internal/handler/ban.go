package handler

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *Handler) HandleBan(
	ctx context.Context,
	event eventsub.ChannelBanEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel ban",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
		slog.String("moderatorName", event.ModeratorUserName),
		slog.String("moderatorId", event.ModeratorUserID),
	)

	go func() {
		channel := deprecatedmodel.Channels{}
		if err := c.gorm.
			WithContext(ctx).
			Where(`"id" = ?`, event.BroadcasterUserID).
			First(&channel).
			Error; err != nil {
			c.logger.Error("channel not found", slog.Any("err", err))
			return
		}

		if channel.BotID == event.UserId {
			channel.IsEnabled = false
			if err := c.gorm.WithContext(ctx).Save(&channel).Error; err != nil {
				c.logger.Error("failed to disable channel", slog.Any("err", err))
			}

			return
		}
	}()

	banEndsIn := event.EndsAt.Sub(time.Now().UTC())
	endsAt := lo.If(event.IsPermanent, "permanent").Else(
		fmt.Sprintf(
			"%v",
			math.Round(banEndsIn.Minutes()),
		),
	)

	c.twirBus.Events.ChannelBan.Publish(
		ctx,
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
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserId,
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

	c.websocketsGrpc.DudesUserPunished(
		ctx,
		&websockets.DudesUserPunishedRequest{
			ChannelId:       event.BroadcasterUserID,
			UserId:          event.UserId,
			UserDisplayName: event.UserName,
			UserName:        event.UserLogin,
		},
	)
}
