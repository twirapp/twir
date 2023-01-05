package command_counter

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"strconv"

	model "github.com/satont/tsuwari/libs/gomodels"

	"gorm.io/gorm"
)

func getCount(commandId string, userId *string) (string, error) {
	db := do.MustInvoke[gorm.DB](di.Provider)

	var count int64

	tx := db.Model(&model.ChannelsCommandsUsages{}).Where(`"commandId" = ?`, commandId)

	if userId != nil {
		tx = tx.Where(`"commandId" = ? AND "userId" = ?`, commandId, *userId)
	}

	err := tx.Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	s := strconv.FormatInt(count, 10)

	return s, nil
}
