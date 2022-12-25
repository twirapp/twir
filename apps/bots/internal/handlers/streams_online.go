package handlers

import (
	"encoding/json"
	"fmt"

	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

type streamOnlineData struct {
	StreamID  string `json:"streamId"`
	ChannelID string `json:"channelId"`
}

func StreamsOnline(db *gorm.DB, data string) {
	streamOnlineStruct := &streamOnlineData{}
	if err := json.Unmarshal([]byte(data), &streamOnlineStruct); err != nil {
		fmt.Println(err)
		return
	}

	channel := model.Channels{}
	if err := db.Where("id = ?", streamOnlineStruct.ChannelID).Find(&channel).Error; err != nil {
		fmt.Println(err)
		return
	}

	if channel.ID == "" {
		return
	}

	err := db.Model(&model.ChannelsGreetings{}).
		Where(`"channelId" = ?`, channel.ID).
		Update("processed", false).Error
	if err != nil {
		fmt.Println(err)
	}
}
