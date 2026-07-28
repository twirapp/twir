package manager

import (
	"sync"
	"time"

	"github.com/google/uuid"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
)

type TimerID uuid.UUID

func (t TimerID) String() string {
	return uuid.UUID(t).String()
}

type Timer struct {
	tickMu sync.Mutex

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

func (t *Timer) withTickLock(fn func()) {
	t.tickMu.Lock()
	defer t.tickMu.Unlock()

	fn()
}
