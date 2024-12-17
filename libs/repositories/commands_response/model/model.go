package model

import (
	"github.com/google/uuid"
)

type Response struct {
	ID                uuid.UUID
	Text              *string
	CommandID         uuid.UUID
	Order             int
	TwitchCategoryIDs []string `db:"twitch_category_id"`
}

var Nil = Response{}
