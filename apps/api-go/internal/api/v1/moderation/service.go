package moderation

import (
	"fmt"
	model "tsuwari/models"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
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
	dto []moderationDto,
	services types.Services,
) ([]model.ChannelsModerationSettings, error) {
	fmt.Printf("%+v\n", dto)
	return nil, nil
}
