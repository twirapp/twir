package scheduledvipsentity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RemoveType string

const (
	RemoveTypeTime      RemoveType = "time"
	RemoveTypeStreamEnd RemoveType = "stream_end"
)

type ScheduledVip struct {
	ID         ulid.ULID
	UserID     string
	ChannelID  string
	CreatedAt  time.Time
	RemoveAt   *time.Time
	RemoveType *RemoveType

	isNil bool
}

func (c ScheduledVip) IsNil() bool {
	return c.isNil
}

var Nil = ScheduledVip{isNil: true}
