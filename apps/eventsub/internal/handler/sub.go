package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
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

func (c *Handler) HandleChannelSubscribe(
	ctx context.Context,
	event eventsub.ChannelSubscribeEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	level := getSubPlan(string(event.Tier))

	c.logger.Info(
		"channel subscribe",
		slog.String("channel_id", event.BroadcasterUserId),
		slog.String("user_id", event.UserId),
		slog.String("level", level),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
			Type:      model.ChannelEventListItemTypeSubscribe,
			Data: &model.ChannelsEventsListItemData{
				SubUserName:        event.UserLogin,
				SubUserDisplayName: event.UserName,
				SubLevel:           level,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	if err := c.twirBus.Events.Subscribe.Publish(
		ctx,
		events.SubscribeMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:          event.UserId,
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Level:           level,
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}

// resub
func (c *Handler) HandleChannelSubscriptionMessage(
	ctx context.Context,
	event eventsub.ChannelSubscriptionMessageEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	level := getSubPlan(string(event.Tier))
	c.logger.Info(
		"channel resubscribe",
		slog.String("channel_id", event.BroadcasterUserId),
		slog.String("user_id", event.UserId),
		slog.String("level", level),
		slog.Int("months", event.CumulativeTotal),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
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
		c.logger.Error(err.Error(), logger.Error(err))
	}

	if err := c.twirBus.Events.ReSubscribe.Publish(
		ctx,
		events.ReSubscribeMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserID:          event.UserId,
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			Months:          int64(event.CumulativeTotal),
			Streak:          int64(event.StreakMonths),
			IsPrime:         level == "prime",
			Message:         event.Message.Text,
			Level:           level,
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}
