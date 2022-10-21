package greetings

import (
	model "tsuwari/models"

	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func HandleGet(channelId string, services types.Services) []model.ChannelsGreetings {
	greetings := []model.ChannelsGreetings{}
	services.DB.Where(`"channelId" = ?`, channelId).Find(&greetings)

	return greetings
}
