package grpc_impl

import (
	"time"
)

const expireShift = 5 * time.Second

func isTokenExpired(expiresIn int, obtainmentTimestamp time.Time) bool {
	currentTime := time.Now().UTC()
	currentTokenLiveTime := currentTime.Sub(obtainmentTimestamp.UTC())

	return int64(currentTokenLiveTime.Seconds())+int64(expireShift.Seconds()) >= int64(expiresIn)
}
