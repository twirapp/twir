package model

import (
	"time"
)

type Stream struct {
	ID           string
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

	isNil bool
}

func (c Stream) IsNil() bool {
	return c.isNil || c.ID == ""
}

var Nil = Stream{
	isNil: true,
}
