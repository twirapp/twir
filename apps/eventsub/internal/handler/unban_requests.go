package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *Handler) HandleChannelUnbanRequestCreate(
	ctx context.Context,
	event eventsub.ChannelUnbanRequestCreateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel unban request create",
		slog.Group(
			"channel",
			slog.String("id", event.BroadcasterUserId),
			slog.String("login", event.BroadcasterUserLogin),
		),
		slog.Group("user", slog.String("id", event.UserId), slog.String("login", event.UserLogin)),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
			Type:      model.ChannelEventListItemTypeChannelUnbanRequestCreate,
			Data: &model.ChannelsEventsListItemData{
				UserLogin:       event.UserLogin,
				UserDisplayName: event.UserName,
				Message:         event.Text,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.twirBus.Events.ChannelUnbanRequestCreate.Publish(
		ctx,
		events.ChannelUnbanRequestCreateMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:             event.UserName,
			UserLogin:            event.UserLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			Text:                 event.Text,
		},
	)
}

func (c *Handler) HandleChannelUnbanRequestResolve(
	ctx context.Context,
	event eventsub.ChannelUnbanRequestResolveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel unban request resolve",
		slog.Group(
			"channel",
			slog.String("id", event.BroadcasterUserId),
			slog.String("login", event.BroadcasterUserLogin),
		),
		slog.Group("user", slog.String("id", event.UserId), slog.String("login", event.UserLogin)),
		slog.Group(
			"moderator",
			slog.String("id", event.ModeratorUserId),
			slog.String("login", event.ModeratorUserLogin),
		),
		slog.String("status", string(event.Status)),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
			Type:      model.ChannelEventListItemTypeChannelUnbanRequestResolve,
			Data: &model.ChannelsEventsListItemData{
				UserLogin:            event.UserLogin,
				UserDisplayName:      event.UserName,
				ModeratorName:        event.ModeratorUserName,
				ModeratorDisplayName: event.ModeratorUserLogin,
				Message:              event.ResolutionText,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.twirBus.Events.ChannelUnbanRequestResolve.Publish(
		ctx,
		events.ChannelUnbanRequestResolveMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:             event.UserName,
			UserLogin:            event.UserLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ModeratorUserName:    event.ModeratorUserName,
			ModeratorUserLogin:   event.ModeratorUserLogin,
			Declined:             event.Status != eventsub.ChannelUnbanRequestResolveEventStatusApproved,
			Reason:               event.ResolutionText,
		},
	)
}
