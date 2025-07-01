package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelUnbanRequestCreate(
	ctx context.Context,
	_ *esb.ResponseHeaders,
	event *esb.ChannelUnbanRequestCreate,
) {
	c.logger.Info(
		"channel unban request create",
		slog.Group(
			"channel",
			slog.String("id", event.BroadcasterUserID),
			slog.String("login", event.BroadcasterUserLogin),
		),
		slog.Group("user", slog.String("id", event.UserID), slog.String("login", event.UserLogin)),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
			Type:      model.ChannelEventListItemTypeChannelUnbanRequestCreate,
			Data: &model.ChannelsEventsListItemData{
				UserLogin:       event.UserLogin,
				UserDisplayName: event.UserName,
				Message:         event.Text,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.twirBus.Events.ChannelUnbanRequestCreate.Publish(
		ctx,
		events.ChannelUnbanRequestCreateMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
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

func (c *Handler) handleChannelUnbanRequestResolve(
	ctx context.Context,
	r *esb.ResponseHeaders,
	event *esb.ChannelUnbanRequestResolve,
) {
	c.logger.Info(
		"channel unban request resolve",
		slog.Group(
			"channel",
			slog.String("id", event.BroadcasterUserID),
			slog.String("login", event.BroadcasterUserLogin),
		),
		slog.Group("user", slog.String("id", event.UserID), slog.String("login", event.UserLogin)),
		slog.Group(
			"moderator",
			slog.String("id", event.ModeratorUserID),
			slog.String("login", event.ModeratorUserLogin),
		),
		slog.String("status", string(event.Status)),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.twirBus.Events.ChannelUnbanRequestResolve.Publish(
		ctx,
		events.ChannelUnbanRequestResolveMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:             event.UserName,
			UserLogin:            event.UserLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ModeratorUserName:    event.ModeratorUserName,
			ModeratorUserLogin:   event.ModeratorUserLogin,
			Declined:             event.Status != esb.ChannelUnbanRequestResolveStatusApproved,
			Reason:               event.ResolutionText,
		},
	)
}
