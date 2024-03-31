package model

import (
	"github.com/satont/twir/libs/types/types/events"
)

type EventOperationFilter struct {
	ID          string                          `gorm:"primaryKey;column:id;type:UUID;default:uuid_generate_v4()" json:"id"`
	OperationID string                          `gorm:"column:operationId;type:UUID;" json:"operationId"`
	Type        events.EventOperationFilterType `gorm:"column:type;type:TEXT;" json:"type"`
	Left        string                          `gorm:"column:left;type:TEXT;" json:"left"`
	Right       string                          `gorm:"column:right;type:TEXT;" json:"right"`

	Operation *EventOperation `gorm:"foreignkey:OperationID" json:"operation"`
}

func (c *EventOperationFilter) TableName() string {
	return "channels_events_operations_filters"
}
