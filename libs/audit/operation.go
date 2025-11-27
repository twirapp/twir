package audit

type OperationMetadata struct {
	System    string
	ActorID   *string
	ChannelID *string
	ObjectID  *string
}

type CreateOperation struct {
	Metadata OperationMetadata
	NewValue any
}

type DeleteOperation struct {
	Metadata OperationMetadata
	OldValue any
}

type UpdateOperation struct {
	Metadata OperationMetadata
	NewValue any
	OldValue any
}
