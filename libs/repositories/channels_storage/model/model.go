package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ChannelStorage struct {
	ID        uuid.UUID       `db:"id"`
	ChannelID string          `db:"channel_id"`
	Key       string          `db:"key"`
	Value     json.RawMessage `db:"value"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`

	isNil bool
}

func (c ChannelStorage) IsNil() bool {
	return c.isNil
}

var Nil = ChannelStorage{
	isNil: true,
}
