package model

import (
	"github.com/google/uuid"
)

type EventsubTopic struct {
	ID            uuid.UUID             `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Topic         string                `gorm:"column:topic;type:varchar(255);not null"`
	Version       string                `gorm:"column:version;type:varchar(255);not null"`
	ConditionType EventsubConditionType `gorm:"column:condition_type;type:int;not null"`

	Subscriptions []EventsubSubscription `gorm:"foreignKey:TopicID;references:ID"`
}

func (EventsubTopic) TableName() string {
	return "eventsub_topics"
}

type EventsubConditionType string

const (
	// EventsubConditionTypeBroadcasterUserID { "broadcaster_user_id": "1234" }
	EventsubConditionTypeBroadcasterUserID EventsubConditionType = "BROADCASTER_USER_ID"
	// EventsubConditionTypeUserID { "user_id": "1234" }
	EventsubConditionTypeUserID EventsubConditionType = "USER_ID"
	// EventsubConditionTypeBroadcasterWithUserID { "broadcaster_user_id": "1234", "user_id": "1234" }
	EventsubConditionTypeBroadcasterWithUserID EventsubConditionType = "BROADCASTER_WITH_USER_ID"
	// EventsubConditionTypeBroadcasterWithModeratorID { "broadcaster_user_id": "1234", "moderator_user_id": "1234" }
	EventsubConditionTypeBroadcasterWithModeratorID EventsubConditionType = "BROADCASTER_WITH_MODERATOR_ID"
	// EventsubConditionTypeToBroadcasterID { "to_broadcaster_user_id": userId }
	EventsubConditionTypeToBroadcasterID EventsubConditionType = "TO_BROADCASTER_ID"
)
