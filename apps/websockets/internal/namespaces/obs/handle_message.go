package obs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/websockets/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"time"
)

func (c *OBS) handleMessage(session *melody.Session, msg []byte) {
	userId, ok := session.Get("userId")
	if userId == "" || !ok {
		return
	}

	data := &types.WebSocketMessage{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	if data.EventName == "setSources" {
		bytes, _ := json.Marshal(data.Data)
		var scenesData map[string][]obsSource
		err = json.Unmarshal(bytes, &scenesData)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}
		c.handleSetSources(userId.(string), scenesData)
	}

	if data.EventName == "setAudioSources" {
		bytes, _ := json.Marshal(data.Data)
		var audioSources []obsAudioSource
		err = json.Unmarshal(bytes, &audioSources)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}
		c.handleSetAudioSources(userId.(string), audioSources)
	}

	if data.EventName == "requestSettings" {
		c.handleRequestSettings(userId.(string))
	}
}

type obsSource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type obsAudioSource string

func (c *OBS) handleSetAudioSources(channelId string, sources []obsAudioSource) {
	bytes, _ := json.Marshal(sources)
	err := c.services.Redis.Set(
		context.Background(),
		fmt.Sprintf("obs:audio-sources:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
}

func (c *OBS) handleSetSources(channelId string, scenes map[string][]obsSource) {
	scenesNames := lo.Keys(scenes)
	bytes, _ := json.Marshal(scenesNames)
	err := c.services.Redis.Set(
		context.Background(),
		fmt.Sprintf("obs:scenes:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	var sourceNames []string
	for _, scene := range scenes {
		for _, source := range scene {
			sourceNames = append(sourceNames, source.Name)
		}
	}
	bytes, _ = json.Marshal(sourceNames)
	err = c.services.Redis.Set(
		context.Background(),
		fmt.Sprintf("obs:sources:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
}

func (c *OBS) handleRequestSettings(channelId string) {
	settings := &model.ChannelModulesSettings{}
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").
		Find(settings).
		Error

	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	obsSettings := &modules.OBSWebSocketSettings{}
	err = json.Unmarshal(settings.Settings, obsSettings)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	outCome := &types.WebSocketMessage{
		EventName: "settings",
		Data:      obsSettings,
	}

	bytes, err := json.Marshal(outCome)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	c.manager.BroadcastFilter(bytes, func(session *melody.Session) bool {
		id, ok := session.Get("userId")
		return ok && id.(string) == channelId
	})
}