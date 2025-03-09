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
}

var Nil = Stream{}
