package bot

import (
	"context"
	"net/http"
	"sync"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"

	apiTypes "github.com/satont/tsuwari/libs/types/types/api/bot"
)

func handleGet(channelId string, services types.Services) (*apiTypes.BotInfo, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	client, err := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           services.DB,
		ClientId:     services.Cfg.TwitchClientId,
		ClientSecret: services.Cfg.TwitchClientSecret,
	}).Create(channelId)
	if client == nil || err != nil {
		logger.Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"cannot create twitch client from your tokens. Please try to reauthorize",
		)
	}

	channel := &model.Channels{}
	err = services.DB.Where("id = ?", channelId).First(channel).Error
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
		mods, err := client.GetChannelMods(&helix.GetChannelModsParams{
			BroadcasterID: channelId,
			UserID:        channel.BotID,
		})
		if err != nil {
			logger.Error(err)
			return
		}

		result.IsMod = lo.If(len(mods.Data.Mods) == 0, false).Else(true)
	}()

	go func() {
		defer wg.Done()
		infoReq, err := client.GetUsers(&helix.UsersParams{
			IDs: []string{channel.BotID},
		})
		if err != nil {
			logger.Error(err)
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

func handlePatch(channelId string, dto *connectionDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)
	grpc := do.MustInvoke[bots.BotsClient](di.Injector)

	twitchUsers, err := services.Twitch.Client.GetUsers(
		&helix.UsersParams{IDs: []string{channelId}},
	)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get twitch user")
	}

	if len(twitchUsers.Data.Users) == 0 {
		return fiber.NewError(http.StatusNotFound, "user not found on twitch")
	}

	user := twitchUsers.Data.Users[0]

	dbUser := model.Channels{}
	err = services.DB.Where(`"id" = ?`, channelId).First(&dbUser).Error

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get user from database")
	}

	if dto.Action == "part" {
		dbUser.IsEnabled = false
	} else {
		dbUser.IsEnabled = true
	}

	services.DB.Where(`"id" = ?`, channelId).Select("*").Updates(&dbUser)

	if dbUser.IsEnabled {
		grpc.Join(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	} else {
		grpc.Leave(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	}

	return nil
}
