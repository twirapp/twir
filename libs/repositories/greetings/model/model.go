package model

import (
	"github.com/google/uuid"
)

type Greeting struct {
	ID           uuid.UUID
	ChannelID    string
	UserID       string
	Enabled      bool
	Text         string
	IsReply      bool
	Processed    bool
	WithShoutOut bool `db:"with_shoutout"`
}

var GreetingNil = Greeting{}
