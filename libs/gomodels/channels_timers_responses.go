package model

type ChannelsTimersResponses struct {
	ID         string          `gorm:"primary_key;column:id;type:TEXT;" json:"-"`
	Text       string          `gorm:"column:text;type:TEXT;"           json:"text"`
	IsAnnounce bool            `gorm:"column:isAnnounce;type:BOOL;"     json:"isAnnounce"`
	TimerID    string          `gorm:"column:timerId;type:uuid;"        json:"-"`
	Timer      *ChannelsTimers `gorm:"foreignKey:ID"                    json:"-"`
}

func (ChannelsTimersResponses) TableName() string {
	return "channels_timers_responses"
}
