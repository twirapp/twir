package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID
	Platform        string
	ChannelID       string
	ChannelLogin    string
	ChannelName     string
	UserID          string
	UserName        string
	UserDisplayName string
	UserColor       string
	Text            string
	CreatedAt       time.Time

	MessageID     string                `json:"messageId,omitempty"`
	MessageType   string                `json:"messageType,omitempty"`
	SenderID      string                `json:"senderId,omitempty"`
	AnnounceColor string                `json:"announceColor,omitempty"`
	Badges        []ChatMessageBadge    `json:"badges,omitempty"`
	Fragments     []ChatMessageFragment `json:"fragments,omitempty"`
}

type ChatMessageBadge struct {
	SetID     string `json:"setId"`
	VersionID string `json:"versionId,omitempty"`
	Text      string `json:"text,omitempty"`
}

type ChatMessageFragment struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	EmoteID  string `json:"emoteId,omitempty"`
	EmoteURL string `json:"emoteUrl,omitempty"`
}

type ChatOverlayModerationEventType string

const (
	ChatOverlayModerationEventUserBanned     ChatOverlayModerationEventType = "USER_BANNED"
	ChatOverlayModerationEventMessageDeleted ChatOverlayModerationEventType = "MESSAGE_DELETED"
	ChatOverlayModerationEventChatCleared    ChatOverlayModerationEventType = "CHAT_CLEARED"
)

type ChatOverlayModerationEvent struct {
	Type      ChatOverlayModerationEventType `json:"type"`
	Platform  string                         `json:"platform"`
	UserLogin string                         `json:"userLogin,omitempty"`
	MessageID string                         `json:"messageId,omitempty"`
}
