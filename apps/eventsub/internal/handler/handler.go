package handler

import (
	eventsub_framework "github.com/dnsge/twitch-eventsub-framework"
	"github.com/satont/twir/apps/eventsub/internal/types"
)

type Handler struct {
	Manager  *eventsub_framework.SubHandler
	services *types.Services
}

func NewHandler(services *types.Services) *Handler {
	manager := eventsub_framework.NewSubHandler(true, []byte(services.Config.TwitchClientSecret))

	myHandler := &Handler{
		Manager:  manager,
		services: services,
	}

	manager.HandleChannelUpdate = myHandler.handleChannelUpdate
	manager.HandleStreamOnline = myHandler.handleStreamOnline
	manager.HandleStreamOffline = myHandler.handleStreamOffline
	manager.HandleUserUpdate = myHandler.handleUserUpdate
	manager.HandleChannelFollow = myHandler.handleChannelFollow
	manager.HandleChannelModeratorAdd = myHandler.handleChannelModeratorAdd
	manager.HandleChannelModeratorRemove = myHandler.handleChannelModeratorRemove
	manager.HandleChannelPointsRewardRedemptionAdd = myHandler.handleChannelPointsRewardRedemptionAdd
	manager.HandleChannelPointsRewardRedemptionUpdate = myHandler.handleChannelPointsRewardRedemptionUpdate
	manager.HandleChannelPollBegin = myHandler.handleChannelPollBegin
	manager.HandleChannelPollProgress = myHandler.handleChannelPollProgress
	manager.HandleChannelPollEnd = myHandler.handleChannelPollEnd
	manager.HandleChannelPredictionBegin = myHandler.handleChannelPredictionBegin
	manager.HandleChannelPredictionProgress = myHandler.handleChannelPredictionProgress
	manager.HandleChannelPredictionLock = myHandler.handleChannelPredictionLock
	manager.HandleChannelPredictionEnd = myHandler.handleChannelPredictionEnd
	manager.HandleChannelBan = myHandler.handleBan
	manager.HandleChannelSubscribe = myHandler.handleChannelSubscribe
	manager.HandleChannelSubscriptionMessage = myHandler.handleChannelSubscriptionMessage
	manager.HandleChannelRaid = myHandler.handleChannelRaid
	manager.HandleChannelChatClear = myHandler.handleChannelChatClear
	manager.HandleChannelChatNotification = myHandler.handleChannelChatNotification
	manager.HandleChannelChatMessage = myHandler.handleChannelChatMessage

	return myHandler
}
