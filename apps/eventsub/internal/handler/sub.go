package handler

import (
	"context"
	"strconv"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
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
