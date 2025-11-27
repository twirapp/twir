package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ChatWallSettings struct {
	ID              ulid.ULID
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
	ID                     ulid.ULID
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
	ID        ulid.ULID
	WallID    ulid.ULID
	UserID    string
	Text      string
	CreatedAt time.Time
}
