package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID
	ChannelID       string
	ChannelLogin    string
	ChannelName     string
	UserID          string
	UserName        string
	UserDisplayName string
	UserColor       string
	Text            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
