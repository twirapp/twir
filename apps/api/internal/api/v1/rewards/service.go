package rewards

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/di"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"
)

var cannotGetRewards = fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")

func handleGet(channelId string) ([]helix.ChannelCustomReward, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewUserClient(context.Background(), channelId, config, tokensGrpc)
	if err != nil {
		return nil, cannotGetRewards
	}

	request, err := twitchClient.GetCustomRewards(
		&helix.GetCustomRewardsParams{
			BroadcasterID:         channelId,
			ID:                    "",
			OnlyManageableRewards: false,
		},
	)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get custom rewards of channel")
	}

	return request.Data.ChannelCustomRewards, nil
}
