package moderation

import (
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	uuid "github.com/satori/go.uuid"
)

/* export declare enum SettingsType {
	links = "links",
	blacklists = "blacklists",
	symbols = "symbols",
	longMessage = "longMessage",
	caps = "caps",
	emotes = "emotes"
} */

var modTypes = []string{"links", "blacklists", "symbols", "longMessage", "caps", "emotes"}

func handleGet(
	channelId string,
	services types.Services,
) ([]model.ChannelsModerationSettings, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	settings := []model.ChannelsModerationSettings{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&settings).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	for _, t := range modTypes {
		_, ok := lo.Find(settings, func(s model.ChannelsModerationSettings) bool {
			return s.Type == t
		})
		if ok {
			continue
		}
		newSetting := model.ChannelsModerationSettings{
			ID:        uuid.NewV4().String(),
			ChannelID: channelId,
			Type:      t,
		}

		err := services.DB.Create(&newSetting).Error
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		settings = append(settings, newSetting)
	}

	return settings, nil
}

func handleUpdate(
	channelId string,
	dto *moderationDto,
	services types.Services,
) ([]model.ChannelsModerationSettings, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	settings := []model.ChannelsModerationSettings{}
	for _, item := range dto.Items {
		setting := model.ChannelsModerationSettings{
			ID:                 item.ID,
			Type:               item.Type,
			ChannelID:          channelId,
			Enabled:            *item.Enabled,
			Subscribers:        *item.Subscribers,
			Vips:               *item.Vips,
			BanTime:            int32(item.BanTime),
			BanMessage:         *item.BanMessage,
			WarningMessage:     *item.WarningMessage,
			CheckClips:         null.BoolFromPtr(item.CheckClips),
			TriggerLength:      null.IntFrom(int64(item.TriggerLength)),
			MaxPercentage:      null.IntFrom(int64(item.MaxPercentage)),
			BlackListSentences: item.BlackListSentences,
		}

		err := services.DB.
			Model(&model.ChannelsModerationSettings{}).
			Select("*").
			Where(`"channelId" = ? AND type = ?`, channelId, item.Type).
			Updates(&setting).Error
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"cannot update moderation settings",
			)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func handlePostTitle(channelId string, dto *postTitleDto, services types.Services) (*postTitleResponse, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewUserClient(channelId, config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	resp, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: channelId,
		Title:         dto.Title,
	})
	if err != nil || resp.StatusCode != 204 {
		logger.Error(err)
		logger.Error(resp.ErrorMessage)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &postTitleResponse{
		Title: dto.Title,
	}, nil
}

func handlePostCategory(channelId string, dto *postCategoryDto, services types.Services) (*postCategoryResponse, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewUserClient(channelId, config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	resp, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: channelId,
		GameID:        dto.CategoryId,
	})
	if err != nil {
		logger.Error(err)
		logger.Error(resp.ErrorMessage)
		return nil, err
	}

	return &postCategoryResponse{
		CategoryId: dto.CategoryId,
	}, nil
}
