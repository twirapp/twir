package keywords

import (
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Keywords) getById(keywordId string) *model.ChannelsKeywords {
	keyword := &model.ChannelsKeywords{}
	err := c.services.Gorm.Where("id = ?", keywordId).First(keyword).Error
	if err != nil {
		return nil
	}

	return keyword
}
