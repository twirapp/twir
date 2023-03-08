package bot

import (
	"context"
	"net/http"
	"sync"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"

	apiTypes "github.com/satont/tsuwari/libs/types/types/api/bot"
)

func (c *Bot) getService(channelId string) (*apiTypes.BotInfo, error) {
	twitchClient, err := twitch.NewUserClient(channelId, *c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create twitch client from your tokens. Please try to reauthorize")
	}

	channel := &model.Channels{}
	err = c.services.Gorm.Where("id = ?", channelId).First(channel).Error
	if err != nil || channel == nil {
		return nil, fiber.NewError(http.StatusNotFound, "cannot find channel in db")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	result := apiTypes.BotInfo{
		Enabled: channel.IsEnabled,
	}

	go func() {
		defer wg.Done()

		if channelId == channel.BotID {
			result.IsMod = true
			return
		}

		mods, err := twitchClient.GetModerators(&helix.GetModeratorsParams{
			BroadcasterID: channelId,
			UserIDs:       []string{channel.BotID},
		})
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		result.IsMod = lo.If(len(mods.Data.Moderators) == 0, false).Else(true)
	}()

	go func() {
		defer wg.Done()
		infoReq, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: []string{channel.BotID},
		})
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		if len(infoReq.Data.Users) == 0 {
			return
		}

		result.BotID = channel.BotID
		result.BotName = infoReq.Data.Users[0].Login
	}()

	wg.Wait()

	if result.BotName == "" {
		return nil, fiber.NewError(404, "cannot find bot info on twitch")
	}

	return &result, nil
}

func (c *Bot) patchService(channelId string, dto *connectionDto) error {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	twitchUsers, err := twitchClient.GetUsers(
		&helix.UsersParams{IDs: []string{channelId}},
	)
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get twitch user")
	}

	if len(twitchUsers.Data.Users) == 0 {
		return fiber.NewError(http.StatusNotFound, "user not found on twitch")
	}

	user := twitchUsers.Data.Users[0]

	dbUser := model.Channels{}
	err = c.services.Gorm.Where(`"id" = ?`, channelId).First(&dbUser).Error

	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get user from database")
	}

	if dto.Action == "part" {
		dbUser.IsEnabled = false
	} else {
		dbUser.IsEnabled = true
	}

	c.services.Gorm.Where(`"id" = ?`, channelId).Select("*").Updates(&dbUser)

	if dbUser.IsEnabled {
		c.services.Grpc.Bots.Join(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	} else {
		c.services.Grpc.Bots.Leave(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	}

	return nil
}
