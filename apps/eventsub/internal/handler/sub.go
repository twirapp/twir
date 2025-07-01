package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
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
	ctx context.Context,
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelSubscribe,
) {
	level := getSubPlan(event.Tier)

	c.logger.Info(
		"channel subscribe",
		slog.String("channel_id", event.BroadcasterUserID),
		slog.String("user_id", event.UserID),
		slog.String("level", level),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
			Type:      model.ChannelEventListItemTypeSubscribe,
			Data: &model.ChannelsEventsListItemData{
				SubUserName:        event.UserLogin,
				SubUserDisplayName: event.UserName,
				SubLevel:           level,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.Subscribe.Publish(
		ctx,
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
	ctx context.Context,
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

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
			Type:      model.ChannelEventListItemTypeReSubscribe,
			Data: &model.ChannelsEventsListItemData{
				ReSubUserName:        event.UserLogin,
				ReSubUserDisplayName: event.UserName,
				ReSubLevel:           level,
				ReSubStreak:          strconv.Itoa(event.StreakMonths),
				ReSubMonths:          strconv.Itoa(event.CumulativeTotal),
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.ReSubscribe.Publish(
		ctx,
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
