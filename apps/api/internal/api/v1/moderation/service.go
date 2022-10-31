package moderation

import (
	"net/http"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
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
	settings := []model.ChannelsModerationSettings{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&settings).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
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
			services.Logger.Sugar().Error(err)
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
			BanMessage:         null.StringFromPtr(item.BanMessage),
			WarningMessage:     null.StringFromPtr(item.WarningMessage),
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
			services.Logger.Sugar().Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"cannot update moderation settings",
			)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}
