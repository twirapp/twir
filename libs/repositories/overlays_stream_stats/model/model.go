package model

import (
	"time"

	"github.com/google/uuid"
)

type StreamStatsOverlay struct {
	ID                   uuid.UUID
	ChannelID            string
	Design               string
	Variant              string
	ViewersEnabled       bool
	ViewersMode          string
	PlatformIconsEnabled bool
	MessagesEnabled      bool
	UptimeEnabled        bool
	SubscribersEnabled   bool
	FollowersEnabled     bool
	ViewersColor         string
	MessagesColor        string
	UptimeColor          string
	SubscribersColor     string
	FollowersColor       string
	CustomHTMLEnabled    bool
	CustomHTML           string
	CustomCSS            string
	CreatedAt            time.Time
	UpdatedAt            time.Time

	isNil bool
}

func (s StreamStatsOverlay) IsNil() bool {
	return s.isNil || s.ID == uuid.Nil
}

var Nil = StreamStatsOverlay{
	isNil: true,
}
