package timers

import "errors"

var (
	ErrNotFound = errors.New("timer not found")
)

type Repository interface {
	GetById(id string) (Timer, error)
	GetAll() ([]Timer, error)
	Update(id string, data Timer) error
	UpdateTriggerMessageNumber(id string, number int) error
}

type TimerResponse struct {
	ID         string
	Text       string
	IsAnnounce bool
}

type Timer struct {
	ID                       string
	Name                     string
	Responses                []TimerResponse
	ChannelID                string
	Interval                 int
	LastTriggerMessageNumber int
	MessageInterval          int
}
