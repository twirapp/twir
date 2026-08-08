package spotify

import (
	"errors"
	"fmt"
)

var (
	ErrNotPremium        = errors.New("spotify: not premium")
	ErrNoActiveDevice    = errors.New("spotify: no active device")
	ErrRestrictedDevice  = errors.New("spotify: restricted device")
	ErrRateLimited       = &RateLimitedError{}
	ErrNotConnected      = errors.New("spotify: not connected")
	ErrTrackNotFound     = errors.New("spotify: track not found")
	ErrInsufficientScope = errors.New("spotify: insufficient scope")
	ErrQueueTimeout      = errors.New("spotify: queue timeout")
)

type RateLimitedError struct {
	RetryAfterSeconds int
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("spotify: rate limited, retry after %d seconds", e.RetryAfterSeconds)
}

func (e *RateLimitedError) Is(target error) bool {
	_, ok := target.(*RateLimitedError)
	return ok
}
