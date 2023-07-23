package model

type EventOperationFilterType string

func (f EventOperationFilterType) String() string {
	return string(f)
}

const (
	EventOperationFilterTypeEquals              EventOperationFilterType = "EQUALS"
	EventOperationFilterTypeNotEquals           EventOperationFilterType = "NOT_EQUALS"
	EventOperationFilterTypeContains            EventOperationFilterType = "CONTAINS"
	EventOperationFilterTypeNotContains         EventOperationFilterType = "NOT_CONTAINS"
	EventOperationFilterTypeStartsWith          EventOperationFilterType = "STARTS_WITH"
	EventOperationFilterTypeEndsWith            EventOperationFilterType = "ENDS_WITH"
	EventOperationFilterTypeGreaterThan         EventOperationFilterType = "GREATER_THAN"
	EventOperationFilterTypeLessThan            EventOperationFilterType = "LESS_THAN"
	EventOperationFilterTypeGreaterThanOrEquals EventOperationFilterType = "GREATER_THAN_OR_EQUALS"
	EventOperationFilterTypeLessThanOrEquals    EventOperationFilterType = "LESS_THAN_OR_EQUALS"
	EventOperationFilterTypeIsEmpty             EventOperationFilterType = "IS_EMPTY"
	EventOperationFilterTypeIsNotEmpty          EventOperationFilterType = "IS_NOT_EMPTY"
)

type EventOperationFilter struct {
	ID          string                   `gorm:"primaryKey;column:id;type:UUID;default:uuid_generate_v4()" json:"id"`
	OperationID string                   `gorm:"column:operationId;type:UUID;" json:"operationId"`
	Type        EventOperationFilterType `gorm:"column:type;type:TEXT;" json:"type"`
	Left        string                   `gorm:"column:left;type:TEXT;" json:"left"`
	Right       string                   `gorm:"column:right;type:TEXT;" json:"right"`

	Operation *EventOperation `gorm:"foreignkey:OperationID" json:"operation"`
}

func (c *EventOperationFilter) TableName() string {
	return "channels_events_operations_filters"
}
