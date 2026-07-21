package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type Stream struct {
	ID           string
	ChannelID    uuid.UUID `db:"channel_id"`
	UserId       string
	UserLogin    string
	UserName     string
	GameId       string
	GameName     string
	CommunityIds []string
	Type         string
	Title        string
	ViewerCount  int
	StartedAt    time.Time
	Language     string
	ThumbnailUrl string
	TagIds       []string
	Tags         []string
	IsMature     bool
	Platform     platform.Platform

	isNil bool
}

func (c Stream) IsNil() bool {
	return c.isNil || c.ID == ""
}

var Nil = Stream{
	isNil: true,
}
