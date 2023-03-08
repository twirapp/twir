package obs_websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	modules "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services *types.Services) (*modules.OBSWebSocketSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.Gorm.Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").First(&settings).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "settings not found")
	}

	data := modules.OBSWebSocketSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &data, nil
}

func handlePost(channelId string, dto *modules.OBSWebSocketSettings, services *types.Services) error {
	
	var existedSettings *model.ChannelModulesSettings
	err := services.Gorm.Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").First(&existedSettings).Error

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
			"type":      "obs_websocket",
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

type OBSWebSocketData struct {
	Sources      []string `json:"sources"`
	AudioSources []string `json:"audioSources"`
	Scenes       []string `json:"scenes"`
}

func handleGetData(channelId string, services *types.Services) (*OBSWebSocketData, error) {
	

	ctx := context.Background()

	sourcesReq := services.Redis.Get(ctx, fmt.Sprintf("obs:sources:%s", channelId)).Val()
	audioSourcesReq := services.Redis.Get(ctx, fmt.Sprintf("obs:audio-sources:%s", channelId)).Val()
	scenesReq := services.Redis.Get(ctx, fmt.Sprintf("obs:scenes:%s", channelId)).Val()

	sources := make([]string, 0)
	audioSources := make([]string, 0)
	scenes := make([]string, 0)

	json.Unmarshal([]byte(sourcesReq), &sources)
	json.Unmarshal([]byte(audioSourcesReq), &audioSources)
	json.Unmarshal([]byte(scenesReq), &scenes)

	return &OBSWebSocketData{
		Sources:      sources,
		AudioSources: audioSources,
		Scenes:       scenes,
	}, nil
}
