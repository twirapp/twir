package model

import (
	"github.com/google/uuid"
)

type ChannelCategoryAliase struct {
	ID         uuid.UUID
	ChannelID  string
	Alias      string
	CategoryID string
}
