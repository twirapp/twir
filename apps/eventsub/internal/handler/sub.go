package handler

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
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

	c.logger.Info(
		"channel subscribe",
		slog.String("channel_id", event.BroadcasterUserID),
		slog.String("user_id", event.UserID),
		slog.String("level", level),
	)

	if err := c.gorm.Create(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.Subscribe.Publish(
		events.SubscribeMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:          event.UserID,
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Level:           level,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

// resub
func (c *Handler) handleChannelSubscriptionMessage(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelSubscriptionMessage,
) {
	level := getSubPlan(event.Tier)
	c.logger.Info(
		"channel resubscribe",
		slog.String("channel_id", event.BroadcasterUserID),
		slog.String("user_id", event.UserID),
		slog.String("level", level),
		slog.Int("months", event.CumulativeTotal),
	)

	if err := c.gorm.Create(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.ReSubscribe.Publish(
		events.ReSubscribeMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:          event.UserID,
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Months:          int64(event.CumulativeTotal),
			Streak:          int64(event.StreakMonths),
			IsPrime:         level == "prime",
			Message:         event.Message.Text,
			Level:           level,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
