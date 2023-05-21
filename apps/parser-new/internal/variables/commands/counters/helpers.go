package command_counters

import (
	"strconv"

	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func getCount(db *gorm.DB, commandId string, userId *string) (string, error) {
	var count int64

	tx := db.Model(&model.ChannelsCommandsUsages{}).Where(`"commandId" = ?`, commandId)

	if userId != nil {
		tx = tx.Where(`"commandId" = ? AND "userId" = ?`, commandId, *userId)
	}

	err := tx.Count(&count).Error
	if err != nil {
		return "", nil
	}

	s := strconv.FormatInt(count, 10)

	return s, nil
}
