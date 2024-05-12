package model

import (
	"github.com/google/uuid"
)

type EventsubTopic struct {
	ID            uuid.UUID             `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Topic         string                `gorm:"type:varchar(255);not null"`
	Version       string                `gorm:"type:varchar(255);not null"`
	ConditionType EventsubConditionType `gorm:"type:int;not null"`

	Subscriptions []EventsubSubscription `gorm:"foreignKey:TopicID;references:ID"`
}

type EventsubConditionType int

const (
	// EventsubConditionTypeBroadcasterUserID { "broadcaster_user_id": "1234" }
	EventsubConditionTypeBroadcasterUserID EventsubConditionType = iota
	// EventsubConditionTypeUserID { "user_id": "1234" }
	EventsubConditionTypeUserID
	// EventsubConditionTypeBroadcasterWithUserID { "broadcaster_user_id": "1234", "user_id": "1234" }
	EventsubConditionTypeBroadcasterWithUserID
	// EventsubConditionTypeBroadcasterWithModeratorID { "broadcaster_user_id": "1234", "moderator_user_id": "1234" }
	EventsubConditionTypeBroadcasterWithModeratorID
)
