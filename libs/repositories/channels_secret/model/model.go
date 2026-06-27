package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelSecret struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Description null.String `db:"description"`
	Value       string      `db:"value"`
	ChannelID   uuid.UUID   `db:"channel_id"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`

	isNil bool
}

func (c ChannelSecret) IsNil() bool {
	return c.isNil
}

var Nil = ChannelSecret{
	isNil: true,
}
