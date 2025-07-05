package model

import (
	"time"
)

type UsersStats struct {
	ID                string    `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id" db:"id"`
	UserID            string    `gorm:"column:userId;type:TEXT;"                        json:"userId" db:"userId"`
	ChannelID         string    `gorm:"column:channelId;type:TEXT;"                     json:"channelId" db:"channelId"`
	Messages          int32     `gorm:"column:messages;type:INT4;default:0;"            json:"messages" db:"messages"`
	Watched           int64     `gorm:"column:watched;type:INT8;default:0;"             json:"watched" db:"watched"`
	UsedChannelPoints int64     `gorm:"column:usedChannelPoints;type:INT8;default:0;"   json:"usedChannelPoints" db:"usedChannelPoints"`
	IsMod             bool      `gorm:"column:is_mod;type:BOOL;default:false;"           json:"isMod" db:"is_mod"`
	IsVip             bool      `gorm:"column:is_vip;type:BOOL;default:false;"           json:"isVip" db:"is_vip"`
	IsSubscriber      bool      `gorm:"column:is_subscriber;type:BOOL;default:false;"    json:"isSubscriber" db:"is_subscriber"`
	Reputation        int64     `gorm:"column:reputation;type:INT8;default:0;"              json:"reputation" db:"reputation"`
	Emotes            int       `gorm:"column:emotes;type:INT4;default:0;"              json:"emotes" db:"emotes"`
	CreatedAt         time.Time `gorm:"column:created_at;type:TIMESTAMPTZ;default:now()" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:TIMESTAMPTZ;default:now()" json:"updatedAt"`
}

func (u *UsersStats) TableName() string {
	return "users_stats"
}
