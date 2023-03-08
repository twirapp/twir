package moderation

import (
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
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

func (c *Moderation) getService(
	channelId string,
) ([]model.ChannelsModerationSettings, error) {
	settings := []model.ChannelsModerationSettings{}
	err := c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&settings).Error
	if err != nil {
		c.services.Logger.Error(err)
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

		err := c.services.Gorm.Create(&newSetting).Error
		if err != nil {
			c.services.Logger.Error(err)
			return nil, err
		}
		settings = append(settings, newSetting)
	}

	return settings, nil
}

func (c *Moderation) postService(
	channelId string,
	dto *moderationDto,
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
			BanMessage:         *item.BanMessage,
			WarningMessage:     *item.WarningMessage,
			CheckClips:         null.BoolFromPtr(item.CheckClips),
			TriggerLength:      null.IntFrom(int64(item.TriggerLength)),
			MaxPercentage:      null.IntFrom(int64(item.MaxPercentage)),
			BlackListSentences: item.BlackListSentences,
		}

		err := c.services.Gorm.
			Model(&model.ChannelsModerationSettings{}).
			Select("*").
			Where(`"channelId" = ? AND type = ?`, channelId, item.Type).
			Updates(&setting).Error
		if err != nil {
			c.services.Logger.Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"cannot update moderation settings",
			)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}
