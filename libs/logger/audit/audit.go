package audit

type OperationType string

const (
	OperationCreate      OperationType = "CREATE"
	OperationUpdate      OperationType = "UPDATE"
	OperationDelete      OperationType = "DELETE"
	OperationTypeUnknown OperationType = "UNKNOWN"
)

type Fields struct {
	System        string
	OperationType OperationType
	OldValue      any
	NewValue      any
	ActorID       *string
	ChannelID     *string
	ObjectID      *string
}

type AuditFieldsContextKey struct{}
