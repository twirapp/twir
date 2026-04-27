package generic

import (
	"time"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

type ChatMessage struct {
	Platform          string                  `json:"platform"`
	ChannelID         string                  `json:"channel_id"`
	UserID            string                  `json:"user_id"`
	PlatformChannelID string                  `json:"platform_channel_id"`
	SenderID          string                  `json:"sender_id"`
	SenderLogin       string                  `json:"sender_login"`
	SenderDisplayName string                  `json:"sender_display_name"`
	MessageID         string                  `json:"message_id"`
	Text              string                  `json:"text"`
	Badges            []ChatMessageBadge      `json:"badges,omitempty"`
	Color             string                  `json:"color"`
	Emotes            []ChatMessageEmote      `json:"emotes,omitempty"`
	EnrichedData      ChatMessageEnrichedData `json:"enriched_data,omitempty"`
}

type ChatMessageEnrichedData struct {
	UsedEmotesWithThirdParty map[string]int
	ChannelCommandPrefix     string
	DbChannel                channelsmodel.Channel
	ChannelStream            *streamsmodel.Stream
	DbUser                   *DbUser
	DbUserChannelStat        *DbUserChannelStat
	IsChatterBroadcaster     bool
	IsChatterModerator       bool
	IsChatterVip             bool
	IsChatterSubscriber      bool
}

type DbUser struct {
	ID                string
	TokenID           *string
	IsBotAdmin        bool
	ApiKey            string
	IsBanned          bool
	HideOnLandingPage bool
	CreatedAt         time.Time
}

type DbUserChannelStat struct {
	ID                uuid.UUID
	UserID            string
	ChannelID         string
	Messages          int32
	Watched           int64
	UsedChannelPoints int64
	IsMod             bool
	IsVip             bool
	IsSubscriber      bool
	Reputation        int64
	Emotes            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ChatMessageBadge struct {
	SetID string `json:"set_id"`
	Text  string `json:"text"`
}

type ChatMessageEmote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
