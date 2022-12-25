package model

type ChannelModerationWarn struct {
	ID        string `gorm:"primary_key;column:id"      json:"id"`
	ChannelID string `gorm:"column:channelId;type:text" json:"channelId"`
	UserID    string `gorm:"column:userId;type:text"    json:"userId"`
	Reason    string `gorm:"column:reason;type:text"    json:"reason"`
}

func (ChannelModerationWarn) TableName() string {
	return "channels_moderation_warnings"
}
