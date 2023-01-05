package messages

import (
	"fmt"
	"strings"
	"time"

	model "github.com/satont/tsuwari/libs/gomodels"

	"gorm.io/gorm"
)

func StoreMessage(
	db *gorm.DB,
	messageId, channelId, userId, userName, text string,
	canBeDeleted bool,
) {
	entity := model.ChannelChatMessage{
		MessageId:    messageId,
		ChannelId:    channelId,
		UserId:       userId,
		UserName:     userName,
		Text:         strings.ToLower(text),
		CanBeDeleted: canBeDeleted,
		CreatedAt:    time.Now().UTC(),
	}

	err := db.Create(&entity).Error
	if err != nil {
		fmt.Println(err)
	}
}
