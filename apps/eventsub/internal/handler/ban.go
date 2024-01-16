package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleBan(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelBan,
) {
	zap.S().Infow(
		"channel ban",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"userId", event.UserID,
		"userName", event.UserLogin,
		"moderatorName", event.ModeratorUserName,
		"moderatorId", event.ModeratorUserID,
	)

	t, _ := time.Parse(time.RFC3339, event.EndsAt)
	banEndsIn := t.Sub(time.Now().UTC())
	endsAt := lo.If(event.IsPermanent, "permanent").Else(
		fmt.Sprintf(
			"%v",
			math.Round(banEndsIn.Minutes()),
		),
	)

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
			EndsAt:               endsAt,
			IsPermanent:          event.IsPermanent,
		},
	)

	c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeChannelBan,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				BanReason:            event.Reason,
				BanEndsInMinutes:     endsAt,
				BannedUserLogin:      event.UserLogin,
				BannedUserName:       event.UserName,
				ModeratorDisplayName: event.ModeratorUserName,
				ModeratorName:        event.ModeratorUserLogin,
			},
		},
	)
}
