package chat_client

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *ChatClient) onNotice(message irc.UserNoticeMessage) {
	if message.Tags["msg-id"] == "subgift" {
		level := getSubPlan(message.MsgParams["msg-param-sub-plan"])

		if err := c.services.DB.Create(
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
			c.services.Logger.Error(
				"cannot create subgift entity",
				slog.String("channelId", message.RoomID),
				slog.String("userId", message.MsgParams["user-id"]),
			)
		}
		c.services.EventsGrpc.SubGift(
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
