package timersentity

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type Timer struct {
	ID                       uuid.UUID
	ChannelID                string
	Name                     string
	Enabled                  bool
	OfflineEnabled           bool
	OnlineEnabled            bool
	TimeInterval             int
	MessageInterval          int
	LastTriggerMessageNumber int
	Responses                []Response

	isNil bool
}

var Nil = Timer{
	isNil: true,
}

type Response struct {
	ID            uuid.UUID
	Text          string
	IsAnnounce    bool
	TimerID       uuid.UUID
	Count         int
	AnnounceColor AnnounceColor

	isNil bool
}

var NilResponse = Response{
	isNil: true,
}

func (c Timer) IsNil() bool {
	return c.isNil
}

func (c Response) IsNil() bool {
	return c.isNil
}

type AnnounceColor int

func (c AnnounceColor) String() string {
	return [...]string{"random", "primary", "blue", "green", "orange", "purple"}[c]
}

func (c AnnounceColor) Scan(value interface{}) error {
	if value == nil {
		c = AnnounceColorPrimary
		return nil
	}
	switch v := value.(type) {
	case int64:
		c = AnnounceColor(v)
	case int32:
		c = AnnounceColor(v)
	case int:
		c = AnnounceColor(v)
	default:
		c = AnnounceColorPrimary
	}
	return nil
}

func (c *AnnounceColor) Value() (driver.Value, error) {
	return int(*c), nil
}

const (
	AnnounceColorRandom                = -1
	AnnounceColorPrimary AnnounceColor = iota - 1
	AnnounceColorBlue
	AnnounceColorGreen
	AnnounceColorOrange
	AnnounceColorPurple
)
