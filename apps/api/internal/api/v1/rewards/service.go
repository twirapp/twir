package rewards

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"net/http"
)

var cannotGetRewards = fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")

func handleGet(channelId string) ([]helix.ChannelCustomReward, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewUserClient(channelId, config, tokensGrpc)

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
