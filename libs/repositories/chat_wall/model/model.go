package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatWallSettings struct {
	ID              uuid.UUID
	ChannelID       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MuteSubscribers bool
	MuteVips        bool

	isNil bool
}

func (c ChatWallSettings) IsNil() bool {
	return c.isNil
}

var ChatWallSettingsNil = ChatWallSettings{
	isNil: true,
}

type ChatWallAction string

const (
	ChatWallActionDelete  ChatWallAction = "DELETE"
	ChatWallActionBan     ChatWallAction = "BAN"
	ChatWallActionTimeout ChatWallAction = "TIMEOUT"
)

type ChatWall struct {
	ID                     uuid.UUID
	ChannelID              string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Phrase                 string
	Enabled                bool
	Action                 ChatWallAction
	DurationSeconds        int
	TimeoutDurationSeconds *int
	AffectedMessages       int

	isNil bool
}

func (c ChatWall) IsNil() bool {
	return c.isNil
}

var ChatWallNil = ChatWall{
	isNil: true,
}

type ChatWallLog struct {
	ID        uuid.UUID
	WallID    uuid.UUID
	UserID    string
	Text      string
	CreatedAt time.Time
}
