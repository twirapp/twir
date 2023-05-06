package model

type ChannelCategoryAlias struct {
	ID        string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:UUID;"  json:"id"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"  json:"channelId"`
	Category  string `gorm:"column:category;type:TEXT;"  json:"category"`
	Alias     string `gorm:"column:alias;type:TEXT;"  json:"alias"`
}

func (c *ChannelCategoryAlias) TableName() string {
	return "channels_categories_aliases"
}
