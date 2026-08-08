package model

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type Response struct {
	ID                uuid.UUID           `db:"id"                 json:"id"`
	Text              *string             `db:"text"               json:"text"`
	CommandID         uuid.UUID           `db:"commandId"          json:"commandId"`
	Order             int                 `db:"order"              json:"order"`
	TwitchCategoryIDs []string            `db:"twitch_category_id" json:"twitch_category_ids"`
	OnlineOnly        bool                `db:"online_only"        json:"online_only"`
	OfflineOnly       bool                `db:"offline_only"       json:"offline_only"`
	Platforms         []platform.Platform `db:"platforms"          json:"platforms"`
}

var Nil = Response{}
