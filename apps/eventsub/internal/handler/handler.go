package handler

import (
	eventsub_framework "github.com/dnsge/twitch-eventsub-framework"
	"github.com/satont/tsuwari/apps/eventsub/internal/types"
)

type handler struct {
	Manager  *eventsub_framework.SubHandler
	services *types.Services
}

func NewHandler(services *types.Services) *handler {
	manager := eventsub_framework.NewSubHandler(true, []byte(services.Config.TwitchClientSecret))

	myHandler := &handler{
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

	return myHandler
}
