package rewards

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/libs/twitch"
)

var cannotGetRewards = fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")

func (c *Rewards) getService(channelId string) ([]helix.ChannelCustomReward, error) {
	twitchClient, err := twitch.NewUserClient(channelId, *c.services.Config, c.services.Grpc.Tokens)
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
