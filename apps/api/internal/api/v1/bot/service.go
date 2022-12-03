package bot

import (
	"context"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func handleGet(channelId string, services types.Services) (*bool, error) {
	client, err := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           services.DB,
		ClientId:     services.Cfg.TwitchClientId,
		ClientSecret: services.Cfg.TwitchClientSecret,
	}).Create(channelId)
	if client == nil || err != nil {
		services.Logger.Sugar().Error(err)
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

	mods, err := client.GetChannelMods(&helix.GetChannelModsParams{
		BroadcasterID: channelId,
		UserID:        channel.BotID,
	})
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get mods of channel")
	}

	return lo.ToPtr(lo.If(len(mods.Data.Mods) == 0, false).Else(true)), nil
}

func handlePatch(channelId string, dto *connectionDto, services types.Services) error {
	twitchUsers, err := services.Twitch.Client.GetUsers(
		&helix.UsersParams{IDs: []string{channelId}},
	)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get twitch user")
	}

	if len(twitchUsers.Data.Users) == 0 {
		return fiber.NewError(http.StatusNotFound, "user not found on twitch")
	}

	user := twitchUsers.Data.Users[0]

	dbUser := model.Channels{}
	err = services.DB.Where(`"id" = ?`, channelId).First(&dbUser).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get user from database")
	}

	if dto.Action == "part" {
		dbUser.IsEnabled = false
	} else {
		dbUser.IsEnabled = true
	}

	services.DB.Where(`"id" = ?`, channelId).Select("*").Updates(&dbUser)

	if dbUser.IsEnabled {
		services.BotsGrpc.Join(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	} else {
		services.BotsGrpc.Leave(context.Background(), &bots.JoinOrLeaveRequest{
			BotId:    dbUser.BotID,
			UserName: user.Login,
		})
	}

	return nil
}
