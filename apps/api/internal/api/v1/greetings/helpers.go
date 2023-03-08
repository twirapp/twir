package greetings

import (
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
)

func getTwitchUserByName(userName string, services *types.Services) *helix.User {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		return nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		Logins: []string{userName},
	})

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	twitchUser := twitchUsers.Data.Users[0]
	return &twitchUser
}

func getTwitchUserById(id string, services *types.Services) *helix.User {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		return nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: []string{id},
	})

	if err != nil || len(twitchUsers.Data.Users) == 0 {
		return nil
	}

	twitchUser := twitchUsers.Data.Users[0]
	return &twitchUser
}

func findGreetingByUser(userId string, channelId string, services *types.Services) *model.ChannelsGreetings {
	greeting := &model.ChannelsGreetings{}
	err := services.Gorm.Where(`"channelId" = ? AND "userId" = ?`, channelId, userId).First(greeting).Error
	if err != nil {
		return nil
	}

	return greeting
}

func findGreetingById(id string, services *types.Services) *model.ChannelsGreetings {
	greeting := model.ChannelsGreetings{}
	err := services.Gorm.Where("id = ?", id).First(&greeting).Error
	if err != nil {
		return nil
	}

	return &greeting
}
