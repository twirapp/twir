package commands

import (
	model "tsuwari/models"

	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func HandleGet(channelId string, services types.Services) []model.ChannelsCommands {
	cmds := []model.ChannelsCommands{}
	services.DB.Preload("Responses").Where(`"channelId" = ?`, channelId).Find(&cmds)

	return cmds
}

func HandlePost(channelId string, services types.Services) model.ChannelsCommands {
	return model.ChannelsCommands{}
}
