package youtube

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/olahol/melody"
	"github.com/satont/tsuwari/apps/websockets/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"time"
)

type playEvent struct {
	ID       string `json:"id"`
	Duration int    `json:"duration"`
}

func (c *YouTube) handleMessage(session *melody.Session, msg []byte) {
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

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		c.services.Logger.Error(err)
		return
	}
	if data.EventName == "play" {
		parsedData := &playEvent{}
		err = json.Unmarshal(bytes, parsedData)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		c.handlePlay(userId.(string), parsedData)
	}

	if data.EventName == "skip" {
		parsedData := []string{}
		err = json.Unmarshal(bytes, &parsedData)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		c.handleSkip(userId.(string), parsedData)
	}

	if data.EventName == "reorder" {
		var parsedData []model.RequestedSong
		err = json.Unmarshal(bytes, &parsedData)

		c.handleNewOrder(userId.(string), parsedData)
	}

	if data.EventName == "pause" {
		fmt.Println("get paused")
	}

}

func (c *YouTube) handleSkip(channelId string, ids []string) {
	spew.Dump(ids)
	err := c.services.Gorm.
		Model(&model.RequestedSong{}).
		Where(`id IN (?) AND "channelId" = ?`, ids, channelId).
		Update(`"deletedAt"`, time.Now()).
		Error
	if err != nil {
		c.services.Logger.Error(err)
	}
}

func (c *YouTube) handleNewOrder(channelId string, songs []model.RequestedSong) {
	var count int64
	err := c.services.Gorm.
		Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelId).
		Count(&count).Error
	if err != nil {
		c.services.Logger.Error(err)
		return
	}

	for i, video := range songs {
		err = c.services.Gorm.
			Model(&model.RequestedSong{}).
			Where(`id = ?`, video.ID).
			Update(`"queuePosition"`, i).
			Error
		if err != nil {
			c.services.Logger.Error(err)
			return
		}
	}
}

func (c *YouTube) handlePlay(userId string, data *playEvent) {
	fmt.Println(userId, data)
}
