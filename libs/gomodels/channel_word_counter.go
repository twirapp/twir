package model

type ChannelWordCounter struct {
	ID        string    `gorm:"column:id;type:uuid" json:"id"`
	ChannelID string    `gorm:"column:channelId;type:text" json:"channelId"`
	Channel   *Channels `json:"channel,omitempty"`
	Phrase    string    `gorm:"column:phrase;type:text" json:"phrase"`
	Counter   int32     `gorm:"column:counter;type:int4" json:"counter"`
	Enabled   bool      `gorm:"column:enabled;type:bool" json:"enabled"`
}

func (ChannelWordCounter) TableName() string {
	return "channels_words_counters"
}
