package model

import (
	"github.com/google/uuid"
)

type ChannelFile struct {
	ID        uuid.UUID
	ChannelID string
	FileName  string
	MimeType  string
	Size      int64
}

var ChannelFileNil = ChannelFile{}
