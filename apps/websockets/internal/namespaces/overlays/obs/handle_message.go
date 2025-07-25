package obs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/websockets/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

func (c *OBS) handleMessage(session *melody.Session, msg []byte) {
	userId, ok := session.Get("userId")
	if userId == "" || !ok {
		return
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	if data.EventName == "setSources" {
		bytes, _ := json.Marshal(data.Data)
		var scenesData map[string][]obsSource
		err = json.Unmarshal(bytes, &scenesData)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}
		c.handleSetSources(userId.(string), scenesData)
	}

	if data.EventName == "setAudioSources" {
		bytes, _ := json.Marshal(data.Data)
		var audioSources []obsAudioSource
		err = json.Unmarshal(bytes, &audioSources)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}
		c.handleSetAudioSources(userId.(string), audioSources)
	}

	if data.EventName == "requestSettings" {
		c.handleRequestSettings(userId.(string))
	}

	if data.EventName == "obsConnected" {
		session.Set("obsConnected", true)
	}
}

type obsSource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type obsAudioSource string

func (c *OBS) handleSetAudioSources(channelId string, sources []obsAudioSource) {
	bytes, _ := json.Marshal(sources)
	err := c.redis.Set(
		context.Background(),
		fmt.Sprintf("obs:audio-sources:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
}

func (c *OBS) handleSetSources(channelId string, scenes map[string][]obsSource) {
	scenesNames := lo.Keys(scenes)
	bytes, _ := json.Marshal(scenesNames)
	err := c.redis.Set(
		context.Background(),
		fmt.Sprintf("obs:scenes:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	var sourceNames []string
	for _, scene := range scenes {
		for _, source := range scene {
			sourceNames = append(sourceNames, source.Name)
		}
	}
	bytes, _ = json.Marshal(sourceNames)
	err = c.redis.Set(
		context.Background(),
		fmt.Sprintf("obs:sources:%s", channelId),
		bytes,
		7*24*time.Hour,
	).Err()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
}

func (c *OBS) handleRequestSettings(channelId string) {
	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, channelId, "obs_websocket").
		Find(settings).
		Error

	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	obsSettings := &modules.OBSWebSocketSettings{}
	err = json.Unmarshal(settings.Settings, obsSettings)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	outCome := &types.WebSocketMessage{
		EventName: "settings",
		Data:      obsSettings,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(outCome)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			id, ok := session.Get("userId")
			return ok && id.(string) == channelId
		},
	)
}
