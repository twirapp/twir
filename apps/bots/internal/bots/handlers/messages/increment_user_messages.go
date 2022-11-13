package messages

import (
	"fmt"

	model "github.com/satont/tsuwari/libs/gomodels"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func IncrementUserMessages(db *gorm.DB, userId, channelId string) {
	user := model.Users{}
	if err := db.Where("id = ?", userId).Preload("Stats").Find(&user).Error; err != nil {
		fmt.Println(err)
		return
	}

	// no user found
	if user.ID == "" {
		user.ID = userId
		user.ApiKey = uuid.NewV4().String()
		user.IsBotAdmin = false
		user.IsTester = false
		user.Stats = createStats(userId, channelId)

		if err := db.Create(&user).Error; err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if user.Stats == nil {
			newStats := createStats(userId, channelId)
			db.Create(newStats)
		} else {
			db.Model(&model.UsersStats{}).Where("id = ?", user.Stats.ID).Update("messages", user.Stats.Messages+1)
		}
	}
}

func createStats(userId, channelId string) *model.UsersStats {
	stats := &model.UsersStats{
		UserID:    userId,
		ChannelID: channelId,
		Messages:  1,
		Watched:   0,
	}
	return stats
}
