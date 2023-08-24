package handler

import (
	"context"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/twir/libs/grpc/generated/events"
)

func (c *Handler) handleBan(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelBan,
) {
	c.services.Grpc.Events.ChannelBan(
		context.TODO(), &events.ChannelBanMessage{
			BaseInfo: &events.BaseInfo{
				ChannelId: event.BroadcasterUserID,
			},
			UserName:             event.UserName,
			UserLogin:            event.UserLogin,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ModeratorUserName:    event.ModeratorUserName,
			ModeratorUserLogin:   event.ModeratorUserLogin,
			Reason:               event.Reason,
			EndsAt:               event.EndsAt,
			IsPermanent:          event.IsPermanent,
		},
	)
}
