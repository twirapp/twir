package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MatchState struct {
	ChannelID         uuid.UUID
	Revision          int64
	ProviderTimestamp int64
	Snapshot          json.RawMessage
	UpdatedAt         time.Time
}

type OutboxAction string

const (
	OutboxActionCreate  OutboxAction = "create"
	OutboxActionResolve OutboxAction = "resolve"
	OutboxActionCancel  OutboxAction = "cancel"
)

type OutboxActionInput struct {
	ChannelID uuid.UUID
	MatchID   int64
	Action    OutboxAction
	Sequence  int64
	Payload   json.RawMessage
}

type ClaimedOutboxAction struct {
	ID        uuid.UUID
	LockToken uuid.UUID
	OutboxActionInput
	Attempts int
}
