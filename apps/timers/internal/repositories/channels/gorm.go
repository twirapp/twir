package channels

import (
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func NewGorm(db *gorm.DB) Repository {
	return &gormRepository{
		db,
	}
}

type gormRepository struct {
	db *gorm.DB
}

func (c *gormRepository) GetById(id string) (Channel, error) {
	channel := model.Channels{}
	result := Channel{}

	if err := c.db.Where("id = ?", id).Find(&channel).Error; err != nil {
		return result, err
	}

	result.ID = channel.ID
	result.Enabled = channel.IsEnabled
	result.IsBotMod = channel.IsBotMod

	return result, nil
}
