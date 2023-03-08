package rewards

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"github.com/satont/tsuwari/libs/twitch"
)

var cannotGetRewards = fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")

func handleGet(channelId string, services *types.Services) ([]helix.ChannelCustomReward, error) {
	twitchClient, err := twitch.NewUserClient(channelId, *services.Config, services.Grpc.Tokens)
	if err != nil {
		return nil, cannotGetRewards
	}

	request, err := twitchClient.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID:         channelId,
		ID:                    "",
		OnlyManageableRewards: false,
	})
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")
	}

	return request.Data.ChannelCustomRewards, nil
}
