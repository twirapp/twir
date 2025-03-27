package model

import (
	"github.com/oklog/ulid/v2"
)

type ChannelCategoryAliase struct {
	ID         ulid.ULID
	ChannelID  string
	Alias      string
	CategoryID string
}
