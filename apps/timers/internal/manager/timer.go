package manager

import (
	"time"

	"github.com/google/uuid"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
)

type TimerID uuid.UUID

func (t TimerID) String() string {
	return uuid.UUID(t).String()
}

type Timer struct {
	id     TimerID
	ticker *time.Ticker
	dbRow  timersentity.Timer

	// system data
	lastTriggerTimestamp     time.Time
	lastTriggerMessageNumber int
	lastTriggerOfflineNumber int
	currentResponseIndex     int
	offlineMessageNumber     int
}
