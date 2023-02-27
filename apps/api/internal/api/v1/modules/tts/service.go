package tts

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	modules "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func handleGet(channelId string, services types.Services) (*modules.TTSSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.DB.Where(`"channelId" = ? AND "type" = ?`, channelId, "tts").First(&settings).Error
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

func handlePost(channelId string, dto *modules.TTSSettings, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	var existedSettings *model.ChannelModulesSettings
	err := services.DB.Where(`"channelId" = ? AND "type" = ?`, channelId, "tts").First(&existedSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	bytes, err := json.Marshal(*dto)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedSettings.ID == "" {
		err = services.DB.Model(&model.ChannelModulesSettings{}).Create(map[string]interface{}{
			"id":        uuid.NewV4().String(),
			"type":      "tts",
			"settings":  bytes,
			"channelId": channelId,
		}).Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	} else {
		err = services.DB.Model(existedSettings).Updates(map[string]interface{}{"settings": bytes}).Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	}
}

func handleGetInfo() map[string]any {
	config := do.MustInvoke[cfg.Config](di.Provider)

	result := map[string]any{}
	req.R().SetSuccessResult(&result).Get(fmt.Sprintf("http://%s/info", config.TTSServiceUrl))

	return result
}
