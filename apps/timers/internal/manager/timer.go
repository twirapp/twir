package manager

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/timers/model"
)

type TimerID uuid.UUID

func (t TimerID) String() string {
	return uuid.UUID(t).String()
}

type Timer struct {
	id     TimerID
	ticker *time.Ticker
	dbRow  model.Timer

	// system data
	lastTriggerTimestamp     time.Time
	lastTriggerMessageNumber int
	currentResponseIndex     int
}
