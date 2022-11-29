package model

type UsersStats struct {
	ID        string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	UserID    string `gorm:"column:userId;type:TEXT;"                        json:"userId"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	Messages  int32  `gorm:"column:messages;type:INT4;default:0;"            json:"messages"`
	Watched   int64  `gorm:"column:watched;type:INT8;default:0;"             json:"watched"`
}

func (u *UsersStats) TableName() string {
	return "users_stats"
}
