package timers

import "errors"

var (
	NotFoundError = errors.New("timer not found")
)

type Repository interface {
	GetById(id string) (Timer, error)
	GetAll() ([]Timer, error)
}

type TimerResponse struct {
	ID         string
	Text       string
	IsAnnounce bool
}

type Timer struct {
	ID        string
	Responses []TimerResponse
	ChannelID string
	Interval  int
}
