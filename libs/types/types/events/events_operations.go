package events

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

type EventOperationFilterType string

func (f EventOperationFilterType) String() string {
	return string(f)
}

var AllEventOperationFilterType = []EventOperationFilterType{
	EventOperationFilterTypeEquals,
	EventOperationFilterTypeNotEquals,
	EventOperationFilterTypeContains,
	EventOperationFilterTypeNotContains,
	EventOperationFilterTypeStartsWith,
	EventOperationFilterTypeEndsWith,
	EventOperationFilterTypeGreaterThan,
	EventOperationFilterTypeLessThan,
	EventOperationFilterTypeGreaterThanOrEquals,
	EventOperationFilterTypeLessThanOrEquals,
	EventOperationFilterTypeIsEmpty,
	EventOperationFilterTypeIsNotEmpty,
}

func (f EventOperationFilterType) TSName() string {
	switch f {
	case EventOperationFilterTypeEquals:
		return EventOperationFilterTypeEquals.String()
	case EventOperationFilterTypeNotEquals:
		return EventOperationFilterTypeNotEquals.String()
	case EventOperationFilterTypeContains:
		return EventOperationFilterTypeContains.String()
	case EventOperationFilterTypeNotContains:
		return EventOperationFilterTypeNotContains.String()
	case EventOperationFilterTypeStartsWith:
		return EventOperationFilterTypeStartsWith.String()
	case EventOperationFilterTypeEndsWith:
		return EventOperationFilterTypeEndsWith.String()
	case EventOperationFilterTypeGreaterThan:
		return EventOperationFilterTypeGreaterThan.String()
	case EventOperationFilterTypeLessThan:
		return EventOperationFilterTypeLessThan.String()
	case EventOperationFilterTypeGreaterThanOrEquals:
		return EventOperationFilterTypeGreaterThanOrEquals.String()
	case EventOperationFilterTypeLessThanOrEquals:
		return EventOperationFilterTypeLessThanOrEquals.String()
	case EventOperationFilterTypeIsEmpty:
		return EventOperationFilterTypeIsEmpty.String()
	case EventOperationFilterTypeIsNotEmpty:
		return EventOperationFilterTypeIsNotEmpty.String()
	default:
		return ""
	}
}
