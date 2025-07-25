package streams

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

func (c *gormRepository) GetByChannelId(id string) (Stream, error) {
	entity := model.ChannelsStreams{}
	result := Stream{}
	if err := c.db.Where(`"userId" = ?`, id).Find(&entity).Error; err != nil {
		return result, nil
	}

	if entity.ID == "" {
		return result, NotFound
	}

	result.ID = entity.ID
	result.UserLogin = entity.UserLogin
	result.UserID = entity.UserId

	return result, nil
}
