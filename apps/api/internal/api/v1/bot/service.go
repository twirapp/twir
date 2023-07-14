package bot

import (
	"context"
	"net/http"
	"sync"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"

	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"

	"github.com/satont/twir/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/types"

	apiTypes "github.com/satont/twir/libs/types/types/api/bot"
)

func handleGet(ctx context.Context, channelId string, services types.Services) (*apiTypes.BotInfo, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewUserClientWithContext(ctx, channelId, config, tokensGrpc)
	if err != nil {
		return nil, fiber.NewError(
			http.StatusInternalServerError, "cannot create twitch client from your tokens. Please try to reauthorize",
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

		if channelId == channel.BotID {
			result.IsMod = true
			return
		}

		mods, err := twitchClient.GetModerators(
			&helix.GetModeratorsParams{
				BroadcasterID: channelId,
				UserIDs:       []string{channel.BotID},
			},
		)
		if err != nil {
			logger.Error(err)
			return
		}

		result.IsMod = lo.If(len(mods.Data.Moderators) == 0, false).Else(true)
	}()

	go func() {
		defer wg.Done()
		infoReq, err := twitchClient.GetUsers(
			&helix.UsersParams{
				IDs: []string{channel.BotID},
			},
		)
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
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	grpc := do.MustInvoke[bots.BotsClient](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	twitchUsers, err := twitchClient.GetUsers(
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
		grpc.Join(
			context.Background(), &bots.JoinOrLeaveRequest{
				BotId:    dbUser.BotID,
				UserName: user.Login,
			},
		)
	} else {
		grpc.Leave(
			context.Background(), &bots.JoinOrLeaveRequest{
				BotId:    dbUser.BotID,
				UserName: user.Login,
			},
		)
	}

	return nil
}
