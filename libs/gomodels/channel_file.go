package model

type ChannelFile struct {
	ID        string `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	ChannelID string `gorm:"column:channel_id;type:text"              json:"channelId"`
	MimeType  string `gorm:"column:mime_type;type:text"              json:"mime_type"`
	Name      string `gorm:"column:file_name;type:text"              json:"file_name"`
	Size      int    `gorm:"column:size;type:INT"              json:"size"`

	Channel *Channels `gorm:"foreignKey:ChannelID" json:"channel"`
}

func (C ChannelFile) TableName() string {
	return "channels_files"
}
