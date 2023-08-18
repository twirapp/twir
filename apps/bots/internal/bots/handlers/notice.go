package handlers

import (
	"context"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"log/slog"
	"strconv"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *Handlers) OnNotice(message irc.UserNoticeMessage) {
	if message.Tags["msg-id"] == "raid" {
		viewers := message.MsgParams["msg-param-viewerCount"]
		intViewers, _ := strconv.Atoi(viewers)

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: message.RoomID,
				UserID:    message.MsgParams["user-id"],
				Type:      model.ChannelEventListItemTypeRaided,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					RaidedViewersCount:    viewers,
					RaidedFromDisplayName: message.MsgParams["msg-param-displayName"],
					RaidedFromUserName:    message.MsgParams["msg-param-login"],
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create raid entity",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
		if _, err := c.eventsGrpc.Raided(
			context.Background(), &events.RaidedMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: message.RoomID},
				UserName:        message.MsgParams["msg-param-login"],
				UserDisplayName: message.MsgParams["msg-param-displayName"],
				Viewers:         int64(intViewers),
				UserId:          message.MsgParams["user-id"],
			},
		); err != nil {
			c.logger.Error(
				"cannot fire raid event",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
	}

	if message.Tags["msg-id"] == "sub" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: message.RoomID,
				Type:      model.ChannelEventListItemTypeSubscribe,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					SubUserName:        message.Tags["login"],
					SubUserDisplayName: message.Tags["display-name"],
					SubLevel:           level,
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create sub entity",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
		if _, err := c.eventsGrpc.Subscribe(
			context.Background(), &events.SubscribeMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: message.RoomID},
				UserName:        message.Tags["login"],
				UserDisplayName: message.Tags["display-name"],
				Level:           level,
				UserId:          message.MsgParams["user-id"],
			},
		); err != nil {
			c.logger.Error(
				"cannot fire sub event",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
	}

	if message.Tags["msg-id"] == "resub" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])
		months, _ := strconv.Atoi(message.MsgParams["msg-param-streak-months"])
		streak, _ := strconv.Atoi(message.MsgParams["msg-param-cumulative-months"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: message.RoomID,
				Type:      model.ChannelEventListItemTypeReSubscribe,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					ReSubUserName:        message.Tags["login"],
					ReSubUserDisplayName: message.Tags["display-name"],
					ReSubLevel:           level,
					ReSubStreak:          message.MsgParams["msg-param-streak-months"],
					ReSubMonths:          message.MsgParams["msg-param-cumulative-months"],
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create resub entity",
				slog.String("channelId", message.RoomID),
				slog.String("userLogin", message.Tags["login"]),
			)
		}
		if _, err := c.eventsGrpc.ReSubscribe(
			context.Background(), &events.ReSubscribeMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: message.RoomID},
				UserName:        message.Tags["login"],
				UserDisplayName: message.Tags["display-name"],
				Months:          int64(months),
				Streak:          int64(streak),
				IsPrime:         level == "prime",
				Message:         message.Message,
				Level:           level,
				UserId:          message.MsgParams["user-id"],
			},
		); err != nil {
			c.logger.Error(
				"cannot fire resub entity",
				slog.String("channelId", message.RoomID),
				slog.String("userLogin", message.Tags["login"]),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
	}

	if message.Tags["msg-id"] == "subgift" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: message.RoomID,
				Type:      model.ChannelEventListItemTypeSubGift,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					SubGiftLevel:                 level,
					SubGiftUserName:              message.Tags["login"],
					SubGiftUserDisplayName:       message.Tags["display-name"],
					SubGiftTargetUserName:        message.MsgParams["msg-param-recipient-user-name"],
					SubGiftTargetUserDisplayName: message.MsgParams["msg-param-recipient-display-name"],
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create subgift entity",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
		c.eventsGrpc.SubGift(
			context.Background(), &events.SubGiftMessage{
				BaseInfo:          &events.BaseInfo{ChannelId: message.RoomID},
				SenderUserName:    message.Tags["login"],
				SenderDisplayName: message.Tags["display-name"],
				TargetUserName:    message.MsgParams["msg-param-recipient-user-name"],
				TargetDisplayName: message.MsgParams["msg-param-recipient-display-name"],
				Level:             level,
				SenderUserId:      message.MsgParams["user-id"],
			},
		)
	}
}

func getSubPlan(plan string) string {
	if plan == "Prime" {
		return "prime"
	}

	parsedPlan, err := strconv.Atoi(plan)
	if err != nil {
		return plan
	}

	return strconv.Itoa(parsedPlan / 1000)
}
