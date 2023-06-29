package messages

import (
	"fmt"

	model "github.com/satont/twir/libs/gomodels"

	"gorm.io/gorm"
)

func IncrementStreamParsedMessages(db *gorm.DB, channelId string) {
	stream := model.ChannelsStreams{}
	if err := db.Where(`"userId" = ?`, channelId).Select("ID", "ParsedMessages").Find(&stream).Error; err != nil {
		fmt.Println(err)
	}
	if stream.ID != "" {
		if err := db.Model(&stream).Update("parsedMessages", stream.ParsedMessages+1).Error; err != nil {
			fmt.Println(err)
		}
	}
}
