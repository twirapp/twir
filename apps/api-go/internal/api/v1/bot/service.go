package bot

import (
	model "tsuwari/models"
	"tsuwari/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func handleGet(channelId string, services types.Services) (*bool, error) {
	client, err := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           services.DB,
		ClientId:     services.Cfg.TwitchClientId,
		ClientSecret: services.Cfg.TwitchClientSecret,
	}).Create(channelId)
	if client == nil || err != nil {
		return nil, fiber.NewError(
			500,
			"cannot create twitch client from your tokens. Please try to reauthorize",
		)
	}

	channel := &model.Channels{}
	err = services.DB.Where("id = ?", channelId).First(channel).Error
	if err != nil || channel == nil {
		return nil, fiber.NewError(404, "cannot find channel in db")
	}

	mods, err := client.GetChannelMods(&helix.GetChannelModsParams{
		BroadcasterID: channelId,
		UserID:        channel.BotID,
	})
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot get mods of channel")
	}

	if len(mods.Data.Mods) == 0 {
		return lo.ToPtr(false), nil
	} else {
		return lo.ToPtr(true), nil
	}
}
