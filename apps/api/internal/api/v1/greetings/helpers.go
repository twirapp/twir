package greetings

import (
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"
)

func getTwitchUserByName(userName string) *helix.User {
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
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

func getTwitchUserById(id string) *helix.User {
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
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

func findGreetingByUser(userId string, channelId string, db *gorm.DB) *model.ChannelsGreetings {
	greeting := &model.ChannelsGreetings{}
	err := db.Where(`"channelId" = ? AND "userId" = ?`, channelId, userId).First(greeting).Error
	if err != nil {
		return nil
	}

	return greeting
}

func findGreetingById(id string, db *gorm.DB) *model.ChannelsGreetings {
	greeting := model.ChannelsGreetings{}
	err := db.Where("id = ?", id).First(&greeting).Error
	if err != nil {
		return nil
	}

	return &greeting
}
