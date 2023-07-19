package handlers

import (
	"context"
	"strconv"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *Handlers) OnNotice(message irc.UserNoticeMessage) {
	if message.Tags["msg-id"] == "raid" {
		viewers := message.MsgParams["msg-param-viewerCount"]
		intViewers, _ := strconv.Atoi(viewers)

		c.eventsGrpc.Raided(
			context.Background(), &events.RaidedMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: message.RoomID},
				UserName:        message.MsgParams["msg-param-login"],
				UserDisplayName: message.MsgParams["msg-param-displayName"],
				Viewers:         int64(intViewers),
				UserId:          message.MsgParams["user-id"],
			},
		)
	}

	if message.Tags["msg-id"] == "sub" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])

		c.eventsGrpc.Subscribe(
			context.Background(), &events.SubscribeMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: message.RoomID},
				UserName:        message.Tags["login"],
				UserDisplayName: message.Tags["display-name"],
				Level:           level,
				UserId:          message.MsgParams["user-id"],
			},
		)
	}

	if message.Tags["msg-id"] == "resub" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])
		months, _ := strconv.Atoi(message.MsgParams["msg-param-streak-months"])
		streak, _ := strconv.Atoi(message.MsgParams["msg-param-cumulative-months"])

		c.eventsGrpc.ReSubscribe(
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
		)
	}

	if message.Tags["msg-id"] == "subgift" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])

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
