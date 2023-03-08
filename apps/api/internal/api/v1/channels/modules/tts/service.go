package tts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	modules "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services *types.Services) (*modules.TTSSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.Gorm.
		Where(`"channelId" = ? AND "type" = ?`, channelId, "tts").
		Where(`"userId" IS NULL`).
		First(&settings).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "settings not found")
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &data, nil
}

func handlePost(channelId string, dto *modules.TTSSettings, services *types.Services) error {
	var existedSettings *model.ChannelModulesSettings
	err := services.Gorm.
		Where(`"channelId" = ? AND "type" = ?`, channelId, "tts").
		Where(`"userId" IS NULL`).
		First(&existedSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	bytes, err := json.Marshal(*dto)
	if err != nil {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedSettings.ID == "" {
		err = services.Gorm.Model(&model.ChannelModulesSettings{}).Create(map[string]interface{}{
			"id":        uuid.NewV4().String(),
			"type":      "tts",
			"settings":  bytes,
			"channelId": channelId,
		}).Error
		if err != nil {
			services.Logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	} else {
		err = services.Gorm.Model(existedSettings).Updates(map[string]interface{}{"settings": bytes}).Error
		if err != nil {
			services.Logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	}
}

func handleGetInfo(services *types.Services) map[string]any {
	result := map[string]any{}
	req.R().SetSuccessResult(&result).Get(fmt.Sprintf("http://%s/info", services.Config.TTSServiceUrl))

	return result
}
