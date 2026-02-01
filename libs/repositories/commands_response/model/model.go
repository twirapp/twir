package model

import (
	"github.com/google/uuid"
)

type Response struct {
	ID                uuid.UUID `db:"id"                 json:"id"`
	Text              *string   `db:"text"               json:"text"`
	CommandID         uuid.UUID `db:"command_id"         json:"command_id"`
	Order             int       `db:"order"              json:"order"`
	TwitchCategoryIDs []string  `db:"twitch_category_id"`
	OnlineOnly        bool      `db:"online_only"        json:"online_only"`
	OfflineOnly       bool      `db:"offline_only"       json:"offline_only"`
}

var Nil = Response{}
