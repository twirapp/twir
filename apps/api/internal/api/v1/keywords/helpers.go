package keywords

import (
	model "github.com/satont/tsuwari/libs/gomodels"

	"gorm.io/gorm"
)

func getById(db *gorm.DB, keywordId string) *model.ChannelsKeywords {
	keyword := &model.ChannelsKeywords{}
	err := db.Where("id = ?", keywordId).First(keyword).Error
	if err != nil {
		return nil
	}

	return keyword
}
