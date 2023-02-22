package model

type ChannelCommandGroup struct {
	ID        string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:UUID;"  json:"id"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"  json:"channelId"`
	Name      string `gorm:"column:name;type:TEXT;"  json:"name"`
	Color     string `gorm:"column:color;type:TEXT;"  json:"color"`

	Commands []ChannelsCommands `gorm:"foreignKey:GroupID" json:"commands"`
}

func (c *ChannelCommandGroup) TableName() string {
	return "channels_commands_groups"
}
