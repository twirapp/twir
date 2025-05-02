package handler

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelUnbanRequestCreate(
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

	c.gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeChannelUnbanRequestCreate,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				UserLogin:       event.UserLogin,
				UserDisplayName: event.UserName,
				Message:         event.Text,
			},
		},
	)

	c.twirBus.Events.ChannelUnbanRequestCreate.Publish(
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

	c.gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeChannelUnbanRequestResolve,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				UserLogin:            event.UserLogin,
				UserDisplayName:      event.UserName,
				ModeratorName:        event.ModeratorUserName,
				ModeratorDisplayName: event.ModeratorUserLogin,
				Message:              event.ResolutionText,
			},
		},
	)

	c.twirBus.Events.ChannelUnbanRequestResolve.Publish(
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
