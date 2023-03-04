package obs_websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	modules "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func handleGet(channelId string, services types.Services) (*modules.OBSWebSocketSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.DB.Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").First(&settings).Error
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

func handlePost(channelId string, dto *modules.OBSWebSocketSettings, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	var existedSettings *model.ChannelModulesSettings
	err := services.DB.Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").First(&existedSettings).Error

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
			"type":      "obs_websocket",
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

type OBSWebSocketData struct {
	Sources      []string `json:"sources"`
	AudioSources []string `json:"audioSources"`
	Scenes       []string `json:"scenes"`
}

func handleGetData(channelId string) (*OBSWebSocketData, error) {
	redisClient := do.MustInvoke[*redis.Client](di.Provider)

	ctx := context.Background()

	sourcesReq := redisClient.Get(ctx, fmt.Sprintf("obs:sources:%s", channelId)).Val()
	audioSourcesReq := redisClient.Get(ctx, fmt.Sprintf("obs:audio-sources:%s", channelId)).Val()
	scenesReq := redisClient.Get(ctx, fmt.Sprintf("obs:scenes:%s", channelId)).Val()

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
