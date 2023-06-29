package messages

import (
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RemoveUserFromLurkers(gorm *gorm.DB, userId string) {
	ignoredUser := &model.IgnoredUser{}
	err := gorm.Where(`"id" = ?`, userId).Find(ignoredUser).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	if ignoredUser.ID != "" {
		gorm.Delete(ignoredUser)
	}
}
