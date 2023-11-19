package handler

import (
	"context"
	"strconv"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/zap"
)

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

func (c *Handler) handleChannelSubscribe(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelSubscribe,
) {
	level := getSubPlan(event.Tier)

	if err := c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			Type:      model.ChannelEventListItemTypeSubscribe,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				SubUserName:        event.UserLogin,
				SubUserDisplayName: event.UserName,
				SubLevel:           level,
			},
		},
	).Error; err != nil {
		zap.S().Error(
			"cannot create sub entity",
			"channelId", event.BroadcasterUserID,
			"userId", event.UserID,
			zap.Error(err),
		)
	}

	if _, err := c.services.Grpc.Events.Subscribe(
		context.Background(), &events.SubscribeMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Level:           level,
			UserId:          event.UserID,
		},
	); err != nil {
		zap.S().Error(
			"cannot fire subscribe event",
			"channelId", event.BroadcasterUserID,
			"userId", event.UserID,
			zap.Error(err),
		)
	}
}

// resub
func (c *Handler) handleChannelSubscriptionMessage(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelSubscriptionMessage,
) {
	level := getSubPlan(event.Tier)

	if err := c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			Type:      model.ChannelEventListItemTypeReSubscribe,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				ReSubUserName:        event.UserLogin,
				ReSubUserDisplayName: event.UserName,
				ReSubLevel:           level,
				ReSubStreak:          strconv.Itoa(event.StreakMonths),
				ReSubMonths:          strconv.Itoa(event.CumulativeTotal),
			},
		},
	).Error; err != nil {
		zap.S().Error(
			"cannot create resub entity",
			"channelId", event.BroadcasterUserID,
			"userLogin", event.UserLogin,
			zap.Error(err),
		)
	}

	if _, err := c.services.Grpc.Events.ReSubscribe(
		context.Background(), &events.ReSubscribeMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Months:          int64(event.CumulativeTotal),
			Streak:          int64(event.StreakMonths),
			IsPrime:         level == "prime",
			Message:         event.Message.Text,
			Level:           level,
			UserId:          event.UserID,
		},
	); err != nil {
		zap.S().Error(
			"cannot fire resub event",
			"channelId", event.BroadcasterUserID,
			"userLogin", event.UserLogin,
			"userId", event.UserID,
			zap.Error(err),
		)
	}
}

// TODO: we should use channel message event when
//
//	https://github.com/dnsge/twitch-eventsub-bindings/issues/2 will be updated and added new
//	bindings
//
// channel message event should be used, not subscription gift,
// because this event doesnt contains target user information
func (c *Handler) handleChannelSubscriptionGift(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelSubscriptionGift,
) {
	// level := getSubPlan(event.Tier)
	// if err := c.services.DB.Create(
	// 	&model.ChannelsEventsListItem{
	// 		ID:        uuid.New().String(),
	// 		ChannelID: message.RoomID,
	// 		Type:      model.ChannelEventListItemTypeSubGift,
	// 		CreatedAt: time.Now().UTC(),
	// 		Data: &model.ChannelsEventsListItemData{
	// 			SubGiftLevel:                 level,
	// 			SubGiftUserName:              message.Tags["login"],
	// 			SubGiftUserDisplayName:       message.Tags["display-name"],
	// 			SubGiftTargetUserName:        message.MsgParams["msg-param-recipient-user-name"],
	// 			SubGiftTargetUserDisplayName: message.MsgParams["msg-param-recipient-display-name"],
	// 		},
	// 	},
	// ).Error; err != nil {
	// 	c.services.Logger.Error(
	// 		"cannot create subgift entity",
	// 		slog.String("channelId", message.RoomID),
	// 		slog.String("userId", message.MsgParams["user-id"]),
	// 	)
	// }
	// c.services.EventsGrpc.SubGift(
	// 	context.Background(), &events.SubGiftMessage{
	// 		BaseInfo:          &events.BaseInfo{ChannelId: message.RoomID},
	// 		SenderUserName:    message.Tags["login"],
	// 		SenderDisplayName: message.Tags["display-name"],
	// 		TargetUserName:    message.MsgParams["msg-param-recipient-user-name"],
	// 		TargetDisplayName: message.MsgParams["msg-param-recipient-display-name"],
	// 		Level:             level,
	// 		SenderUserId:      message.MsgParams["user-id"],
	// 	},
	// )
}
