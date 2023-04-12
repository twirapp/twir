package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/satont/tsuwari/libs/pubsub"

	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func StreamsOffline(db *gorm.DB, data []byte) {
	streamOfflineStruct := &pubsub.StreamOfflineMessage{}
	if err := json.Unmarshal(data, &streamOfflineStruct); err != nil {
		fmt.Println(err)
		return
	}

	channel := model.Channels{}
	if err := db.Where("id = ?", streamOfflineStruct.ChannelID).Find(&channel).Error; err != nil {
		fmt.Println(err)
		return
	}

	if channel.ID == "" {
		return
	}

	//db.Model(&model.ChannelsGreetings{}).
	//	Where(`"channelId" = ?`, channel.ID).
	//	Update("processed", false)
}
