package rewards

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/types"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/twitch"
	"net/http"
)

var cannotGetRewards = fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")

func handleGet(channelId string, services types.Services) ([]helix.ChannelCustomReward, error) {
	config := do.MustInvoke[*cfg.Config](di.Injector)

	api, err := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           services.DB,
		ClientId:     config.TwitchClientId,
		ClientSecret: config.TwitchClientSecret,
	}).Create(channelId)

	if err != nil {
		return nil, cannotGetRewards
	}

	request, err := api.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID:         channelId,
		ID:                    "",
		OnlyManageableRewards: false,
	})

	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")
	}

	return request.Data.ChannelCustomRewards, nil
}
