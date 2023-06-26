package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/satont/twir/libs/pubsub"

	"github.com/satont/twir/apps/bots/internal/bots"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func UserUpdate(db *gorm.DB, botsService *bots.BotsService, data []byte) {
	userStruct := &pubsub.UserUpdateMessage{}
	if err := json.Unmarshal(data, userStruct); err != nil {
		fmt.Println(err)
		return
	}

	channel := model.Channels{}
	if err := db.Where("id = ?", userStruct.UserID).Find(&channel).Error; err != nil {
		fmt.Println(err)
		return
	}

	if channel.ID == "" {
		return
	}

	bot, isBotFound := botsService.Instances[channel.BotID]
	if !isBotFound {
		return
	}

	if channel.IsEnabled {
		bot.Join(userStruct.UserName)
	} else {
		bot.Depart(userStruct.UserName)
	}

	fmt.Printf("%+v\n", data)
}
