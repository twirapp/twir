package entity

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ChannelStorage struct {
	ID        uuid.UUID
	ChannelID string
	Key       string
	Value     json.RawMessage
}

var ChannelStorageNil = ChannelStorage{}
