package entity

import (
	"time"
)

type ChatWallAction string

const (
	ChatWallActionDelete  ChatWallAction = "DELETE"
	ChatWallActionBan     ChatWallAction = "BAN"
	ChatWallActionTimeout ChatWallAction = "TIMEOUT"
)

type ChatWall struct {
	ID                     string
	ChannelID              string
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Phrase                 string
	Enabled                bool
	Action                 ChatWallAction
	DurationSeconds        int
	TimeoutDurationSeconds *int
	AffectedMessages       int
}

type ChatWallLog struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	Text      string
}

type ChatWallSettings struct {
	MuteSubscribers bool
	MuteVips        bool
}
