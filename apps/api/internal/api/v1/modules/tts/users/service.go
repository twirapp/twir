package users

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"go.uber.org/zap"
)

type UserSettings struct {
	Rate  int    `json:"rate"`
	Voice string `json:"voice"`
	Pitch int    `json:"pitch"`

	UserLogin  string `json:"userLogin"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
	UserID     string `json:"userId"`
}

func handleGet(channelId string, services types.Services) ([]*UserSettings, error) {
	var settings []model.ChannelModulesSettings
	err := services.DB.
		Where(`"channelId" = ? AND "type" = ? AND "userId" IS NOT NULL`, channelId, "tts").
		Find(&settings).
		Error
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	var usersSettings []*UserSettings

	for _, setting := range settings {
		var ttsSettings modules.TTSSettings
		err = json.Unmarshal(setting.Settings, &ttsSettings)
		if err != nil {
			zap.S().Error(err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
		}

		usersSettings = append(usersSettings, &UserSettings{
			Rate:   ttsSettings.Pitch,
			Voice:  ttsSettings.Voice,
			Pitch:  ttsSettings.Pitch,
			UserID: setting.UserId.String,
		})
	}

	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	chunks := lo.Chunk(usersSettings, 100)
	wg := &sync.WaitGroup{}
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(c []*UserSettings) {
			defer wg.Done()

			users, err := twitchClient.GetUsers(&helix.UsersParams{
				IDs: lo.Map(c, func(item *UserSettings, _ int) string {
					return item.UserID
				}),
			})

			if err != nil || users.ErrorMessage != "" {
				zap.S().Error(err, users.ErrorMessage)
				return
			}

			for _, user := range users.Data.Users {
				settings, ok := lo.Find(usersSettings, func(s *UserSettings) bool {
					return s.UserID == user.ID
				})
				if !ok {
					continue
				}

				settings.UserAvatar = user.ProfileImageURL
				settings.UserLogin = user.Login
				settings.UserName = user.DisplayName
			}
		}(chunk)
	}

	wg.Wait()

	return usersSettings, nil
}

func handleDelete(channelId, userId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	user := &model.ChannelModulesSettings{}
	err := services.DB.
		Where(`"userId" = ? AND "channelId" = ? AND type = ?`, userId, channelId, "tts").
		Find(user).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	err = services.DB.Delete(user).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func handleDeleteAll(channelId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	err := services.DB.
		Where(`"channelId" = ? AND type = ? and "userId" IS NOT NULL`, channelId, "tts").
		Delete(&model.ChannelModulesSettings{}).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
