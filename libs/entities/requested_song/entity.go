package requested_song

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RequestedSong struct {
	ID                   uuid.UUID
	ChannelID            uuid.UUID
	OrderedByID          uuid.UUID
	OrderedByName        string
	OrderedByDisplayName *string
	VideoID              string
	Title                string
	Duration             int32
	QueuePosition        int
	SongLink             *string
	CreatedAt            time.Time
	DeletedAt            *time.Time

	isNil bool
}

func (c RequestedSong) IsNil() bool {
	return c.isNil
}

var Nil = RequestedSong{isNil: true}

func (c RequestedSong) Link() string {
	if c.SongLink != nil && *c.SongLink != "" {
		return *c.SongLink
	}

	return fmt.Sprintf("https://youtu.be/%s", c.VideoID)
}

func (c RequestedSong) RequesterDisplayName() string {
	if c.OrderedByDisplayName != nil && *c.OrderedByDisplayName != "" {
		return *c.OrderedByDisplayName
	}

	return c.OrderedByName
}
