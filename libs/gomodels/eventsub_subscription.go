package model

import (
	"github.com/google/uuid"
)

type EventsubSubscription struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	TopicID     uuid.UUID `gorm:"column:topic_id;type:uuid;not null"`
	UserID      string    `gorm:"column:user_id;type:text;not null"`
	Status      string    `gorm:"column:status;type:varchar(255);not null"`
	Version     string    `gorm:"column:version;type:varchar(255);not null"`
	CallbackUrl string    `gorm:"column:callback_url;type:varchar(255);not null"`

	User  *Users         `gorm:"foreignKey:UserID;references:ID"`
	Topic *EventsubTopic `gorm:"foreignKey:TopicID;references:ID"`
}
