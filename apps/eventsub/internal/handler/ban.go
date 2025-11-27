package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	deprecatedmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *Handler) handleModerateActionBan(
	ctx context.Context,
	event eventsub.ChannelModerateEventV2,
) {
	var userId, userLogin, userName, reason string
	if event.Ban != nil {
		userId = event.Ban.UserID
		userLogin = event.Ban.UserLogin
		userName = event.Ban.UserName
		reason = event.Ban.Reason
	} else if event.Timeout != nil {
		userId = event.Timeout.UserID
		userLogin = event.Timeout.UserLogin
		userName = event.Timeout.UserName
		reason = event.Timeout.Reason
	}

	go func() {
		channel := deprecatedmodel.Channels{}
		if err := c.gorm.
			WithContext(ctx).
			Where(`"id" = ?`, event.BroadcasterUserID).
			First(&channel).
			Error; err != nil {
			c.logger.Error("channel not found", logger.Error(err))
			return
		}

		if channel.BotID == userId {
			channel.IsEnabled = false
			if err := c.gorm.WithContext(ctx).Save(&channel).Error; err != nil {
				c.logger.Error("failed to disable channel", logger.Error(err))
			}

			return
		}
	}()

	var endsAt string
	if event.Ban != nil {
		endsAt = "permanent"
	} else if event.Timeout != nil {
		banEndsIn := event.Timeout.ExpiresAt.Sub(time.Now().UTC())
		var (
			minutes = banEndsIn.Minutes()
			unit    string
		)
		if minutes == 0 {
			unit = fmt.Sprint(banEndsIn.Seconds())
		} else {
			unit = fmt.Sprint(math.Round(minutes))
		}
		endsAt = unit
	}

	c.twirBus.Events.ChannelBan.Publish(
		ctx,
		events.ChannelBanMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:             userName,
			UserLogin:            userLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserName,
			ModeratorUserName:    event.ModeratorUserName,
			ModeratorUserLogin:   event.ModeratorUserLogin,
			Reason:               reason,
			EndsAt:               endsAt,
			IsPermanent:          event.Ban != nil,
		},
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &userId,
			Type:      model.ChannelEventListItemTypeChannelBan,
			Data: &model.ChannelsEventsListItemData{
				BanReason:            reason,
				BanEndsInMinutes:     endsAt,
				BannedUserLogin:      userLogin,
				BannedUserName:       userName,
				ModeratorDisplayName: event.ModeratorUserName,
				ModeratorName:        event.ModeratorUserLogin,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.websocketsGrpc.DudesUserPunished(
		ctx,
		&websockets.DudesUserPunishedRequest{
			ChannelId:       event.BroadcasterUserID,
			UserId:          userId,
			UserDisplayName: userName,
			UserName:        userLogin,
		},
	)
}

func (c *Handler) handleModerateActionUnBan(
	ctx context.Context,
	event eventsub.ChannelModerateEventV2,
) {
	payload := events.ChannelUnbanMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   event.BroadcasterUserID,
			ChannelName: event.BroadcasterUserLogin,
		},
		BroadcasterUserName:  event.BroadcasterUserName,
		BroadcasterUserLogin: event.BroadcasterUserLogin,
		ModeratorUserID:      event.ModeratorUserID,
		ModeratorUserName:    event.ModeratorUserName,
		ModeratorUserLogin:   event.ModeratorUserLogin,
	}

	if event.Unban != nil {
		payload.UserID = event.Unban.UserID
		payload.UserName = event.Unban.UserName
		payload.UserLogin = event.Unban.UserLogin
	} else if event.Untimeout != nil {
		payload.UserID = event.Untimeout.UserID
		payload.UserName = event.Untimeout.UserName
		payload.UserLogin = event.Untimeout.UserLogin
	}

	c.twirBus.Events.ChannelUnban.Publish(
		ctx,
		payload,
	)
}
