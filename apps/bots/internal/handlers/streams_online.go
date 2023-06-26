package handlers

import (
	"encoding/json"
	"fmt"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/pubsub"
	"gorm.io/gorm"
)

func StreamsOnline(db *gorm.DB, data []byte) {
	streamOnlineStruct := &pubsub.StreamOnlineMessage{}
	if err := json.Unmarshal(data, &streamOnlineStruct); err != nil {
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
