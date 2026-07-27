package model

import (
	"time"

	"github.com/guregu/null"
)

type ChannelsQuotes struct {
	ID          string      `gorm:"primaryKey;column:id;type:UUID" json:"id"`
	ChannelID   string      `gorm:"column:channelId;type:TEXT" json:"channelId"`
	Number      int         `gorm:"column:number;type:INT4" json:"number"`
	Text        string      `gorm:"column:text;type:TEXT" json:"text"`
	CreatorID   null.String `gorm:"column:creatorId;type:TEXT" json:"creatorId"`
	CreatorName null.String `gorm:"column:creatorName;type:TEXT" json:"creatorName"`
	GameID      null.String `gorm:"column:gameId;type:TEXT" json:"gameId"`
	GameName    null.String `gorm:"column:gameName;type:TEXT" json:"gameName"`
	CreatedAt   time.Time   `gorm:"column:createdAt;type:TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"column:updatedAt;type:TIMESTAMP" json:"updatedAt"`
}

func (c *ChannelsQuotes) TableName() string {
	return "channels_quotes"
}
