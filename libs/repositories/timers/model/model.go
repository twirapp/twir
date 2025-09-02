package model

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type Timer struct {
	ID                       uuid.UUID
	ChannelID                string
	Name                     string
	Enabled                  bool
	TimeInterval             int
	MessageInterval          int
	LastTriggerMessageNumber int
	Responses                []Response `db:"responses"`
}

type Response struct {
	ID            uuid.UUID
	Text          string
	IsAnnounce    bool
	TimerID       uuid.UUID
	Count         int
	AnnounceColor AnnounceColor
}

var Nil = Timer{}

type AnnounceColor int

func (c AnnounceColor) String() string {
	return [...]string{"primary", "blue", "green", "orange", "purple"}[c]
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
	AnnounceColorPrimary AnnounceColor = iota
	AnnounceColorBlue
	AnnounceColorGreen
	AnnounceColorOrange
	AnnounceColorPurple
)
