package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/satont/tsuwari/apps/bots/internal/bots"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

type userUpdateUser struct {
	UserID        string `json:"user_id"`
	UserLogin     string `json:"user_login"`
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Description   string `json:"description"`
}

func UserUpdate(db *gorm.DB, botsService *bots.BotsService, data []byte) {
	userStruct := &userUpdateUser{}
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
