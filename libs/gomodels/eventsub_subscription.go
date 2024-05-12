package model

type EventsubSubscription struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TopicID   string `gorm:"type:uuid;not null"`
	ChannelID string `gorm:"type:uuid;not null"`
	Status    string `gorm:"type:varchar(255);not null"`

	Channel *Channels      `gorm:"foreignKey:ChannelID;references:ID"`
	Topic   *EventsubTopic `gorm:"foreignKey:TopicID;references:ID"`
}
