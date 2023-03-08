package greetings

import (
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
)

func (c *Greetings) getTwitchUserByName(userName string) *helix.User {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
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

func (c *Greetings) getTwitchUserById(id string) *helix.User {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
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

func (c *Greetings) findGreetingByUser(userId string, channelId string) *model.ChannelsGreetings {
	greeting := &model.ChannelsGreetings{}
	err := c.services.Gorm.Where(`"channelId" = ? AND "userId" = ?`, channelId, userId).First(greeting).Error
	if err != nil {
		return nil
	}

	return greeting
}

func (c *Greetings) findGreetingById(id string) *model.ChannelsGreetings {
	greeting := model.ChannelsGreetings{}
	err := c.services.Gorm.Where("id = ?", id).First(&greeting).Error
	if err != nil {
		return nil
	}

	return &greeting
}
