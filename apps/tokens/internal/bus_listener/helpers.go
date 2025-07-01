package bus_listener

import (
	"time"
)

const expireShift = 15 * time.Minute

func isTokenExpired(expiresIn int, obtainmentTimestamp time.Time) bool {
	currentTime := time.Now().UTC()
	currentTokenLiveTime := currentTime.Sub(obtainmentTimestamp.UTC())

	return int64(currentTokenLiveTime.Seconds())+int64(expireShift.Seconds()) >= int64(expiresIn)
}
