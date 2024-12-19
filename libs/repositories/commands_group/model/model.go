package model

import (
	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID
	ChannelID string
	Name      string
	Color     string
}

var Nil = Group{}
