package model

import (
	"github.com/google/uuid"
)

type EventsubTopic struct {
	ID      uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Topic   string    `gorm:"column:topic;type:varchar(255);not null"`
	Version string    `gorm:"column:version;type:varchar(255);not null"`

	Subscriptions []EventsubSubscription `gorm:"foreignKey:TopicID;references:ID"`
}

func (EventsubTopic) TableName() string {
	return "eventsub_topics"
}
