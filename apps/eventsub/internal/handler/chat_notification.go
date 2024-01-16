package handler

import (
	"context"
	"time"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelChatNotification(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatNotification,
) {
	switch event.NoticeType {
	case "sub_gift":
		c._notificationSubGift(event)
	case "community_sub_gift":
		c._notificationCommunitySubGift(event)
	}
}

func (c *Handler) _notificationSubGift(event *eventsub_bindings.EventChannelChatNotification) {
	if event.SubGift == nil {
		zap.S().Errorw(
			"subgift event has no subgift data",
			"channelId", event.BroadcasterUserID,
			"channelName", event.BroadcasterUserLogin,
		)
		return
	}

	tier := getSubPlan(event.SubGift.SubTier)

	defer zap.S().Infow(
		"subgift",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"targetUserName", event.SubGift.RecipientUserName,
		"targetUserId", event.SubGift.RecipientUserID,
		"level", event.SubGift.SubTier,
	)

	if err := c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			Type:      model.ChannelEventListItemTypeSubGift,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				SubGiftLevel:                 tier,
				SubGiftUserName:              event.ChatterUserLogin,
				SubGiftUserDisplayName:       event.ChatterUserName,
				SubGiftTargetUserName:        event.SubGift.RecipientUserName,
				SubGiftTargetUserDisplayName: event.SubGift.RecipientUserName,
			},
		},
	).Error; err != nil {
		zap.S().Errorw(
			"cannot create subgift entity",
			"channelId", event.BroadcasterUserID,
			"channelName", event.BroadcasterUserLogin,
			"userId", event.ChatterUserID,
			"userName", event.ChatterUserLogin,
			"error", err,
		)
	}
	c.services.Grpc.Events.SubGift(
		context.Background(), &events.SubGiftMessage{
			BaseInfo:          &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			SenderUserId:      event.ChatterUserID,
			SenderUserName:    event.ChatterUserLogin,
			SenderDisplayName: event.ChatterUserName,
			TargetUserName:    event.SubGift.RecipientUserName,
			TargetDisplayName: event.SubGift.RecipientUserName,
			Level:             tier,
		},
	)
}

func (c *Handler) _notificationCommunitySubGift(
	event *eventsub_bindings.
		EventChannelChatNotification,
) {

}
