package entity

import (
	"time"
)

type SongRequestPublic struct {
	Title           string
	UserID          string
	CreatedAt       time.Time
	SongLink        string
	DurationSeconds int
}
