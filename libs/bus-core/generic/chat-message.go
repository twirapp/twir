package generic

import (
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

type ChatMessage struct {
	Message                     *ChatMessageMessage        `json:"message,omitempty"`
	Cheer                       *ChatMessageCheer          `json:"cheer,omitempty"`
	Reply                       *ChatMessageReply          `json:"reply,omitempty"`
	ID                          string                     `json:"id,omitempty"`
	BroadcasterUserId           string                     `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserName         string                     `json:"broadcaster_user_name,omitempty"`
	BroadcasterUserLogin        string                     `json:"broadcaster_user_login,omitempty"`
	ChatterUserId               string                     `json:"chatter_user_id,omitempty"`
	ChatterUserName             string                     `json:"chatter_user_name,omitempty"`
	ChatterUserLogin            string                     `json:"chatter_user_login,omitempty"`
	MessageType                 string                     `json:"message_type,omitempty"`
	ChannelPointsCustomRewardId string                     `json:"channel_points_custom_reward_id,omitempty"`
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

const (
	broadcasterBadgeId       = "broadcaster"
	subscriberBadgeId        = "subscriber"
	subscriberFounderBadgeId = "founder"
	vipBadgeId               = "vip"
	moderatorBadgeId         = "moderator"
	leadModeratorBadgeId     = "lead_moderator"
)

func (c ChatMessage) IsChatterBroadcaster() bool {
	flag := c.EnrichedData.IsChatterBroadcaster ||
		(c.ChatterUserId != "" && c.BroadcasterUserId != "" && c.ChatterUserId == c.BroadcasterUserId) ||
		(c.SenderID != "" && c.PlatformChannelID != "" && c.SenderID == c.PlatformChannelID) || slices.ContainsFunc(
		c.Badges, func(b ChatMessageBadge) bool {
			return badgeMatchesAny(b, broadcasterBadgeId)
		},
	)

	c.EnrichedData.IsChatterBroadcaster = flag

	return c.EnrichedData.IsChatterBroadcaster
}

func (c ChatMessage) IsChatterVip() bool {
	flag := c.EnrichedData.IsChatterVip || slices.ContainsFunc(
		c.Badges, func(b ChatMessageBadge) bool {
			return badgeMatchesAny(b, vipBadgeId)
		},
	)

	c.EnrichedData.IsChatterVip = flag

	return c.EnrichedData.IsChatterVip
}

func (c ChatMessage) IsChatterSubscriber() bool {
	flag := c.EnrichedData.IsChatterSubscriber || slices.ContainsFunc(
		c.Badges, func(b ChatMessageBadge) bool {
			return badgeMatchesAny(b, subscriberBadgeId, subscriberFounderBadgeId)
		},
	)

	c.EnrichedData.IsChatterSubscriber = flag

	return c.EnrichedData.IsChatterSubscriber
}

func (c ChatMessage) IsChatterModerator() bool {
	flag := c.EnrichedData.IsChatterModerator || slices.ContainsFunc(
		c.Badges, func(b ChatMessageBadge) bool {
			return badgeMatchesAny(b, moderatorBadgeId, leadModeratorBadgeId)
		},
	)

	c.EnrichedData.IsChatterModerator = flag

	return c.EnrichedData.IsChatterModerator
}

func (c ChatMessage) HasRoleFromDbByType(roleType string) bool {
	switch strings.ToLower(roleType) {
	case "broadcaster":
		return c.IsChatterBroadcaster()
	case "moderator":
		return c.IsChatterModerator()
	case "subscriber":
		return c.IsChatterSubscriber()
	case "vip":
		return c.IsChatterVip()
	default:
		return false
	}
}

func badgeMatchesAny(b ChatMessageBadge, values ...string) bool {
	for _, value := range values {
		if strings.EqualFold(b.SetID, value) || strings.EqualFold(b.ID, value) ||
			strings.EqualFold(b.Info, value) || strings.EqualFold(b.Text, value) {
			return true
		}
	}

	return false
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
	ID    string `json:"id,omitempty"`
	SetID string `json:"set_id"`
	Info  string `json:"info,omitempty"`
	Text  string `json:"text"`
}

type ChatMessageEmote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type FragmentType int32

const (
	FragmentType_TEXT      FragmentType = 0
	FragmentType_CHEERMOTE FragmentType = 1
	FragmentType_EMOTE     FragmentType = 2
	FragmentType_MENTION   FragmentType = 3
)

type ChatMessageMessageFragmentPosition struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type ChatMessageMessageFragmentEmote struct {
	ID         string   `json:"id,omitempty"`
	EmoteSetID string   `json:"emote_set_id,omitempty"`
	OwnerID    string   `json:"owner_id,omitempty"`
	Format     []string `json:"format,omitempty"`
}

type ChatMessageMessageFragmentMention struct {
	UserID    string `json:"user_id,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	UserLogin string `json:"user_login,omitempty"`
}

type ChatMessageMessageFragmentCheermote struct {
	Prefix string `json:"prefix,omitempty"`
	Bits   int64  `json:"bits,omitempty"`
	Tier   int64  `json:"tier,omitempty"`
}

type ChatMessageMessageFragment struct {
	Cheermote *ChatMessageMessageFragmentCheermote `json:"cheermote,omitempty"`
	Emote     *ChatMessageMessageFragmentEmote     `json:"emote,omitempty"`
	Mention   *ChatMessageMessageFragmentMention   `json:"mention,omitempty"`
	Text      string                               `json:"text"`
	Position  ChatMessageMessageFragmentPosition   `json:"position,omitempty"`
	Type      FragmentType                         `json:"type"`
}

type ChatMessageMessage struct {
	Text      string                       `json:"text"`
	Fragments []ChatMessageMessageFragment `json:"fragments,omitempty"`
}

type ChatMessageCheer struct {
	Bits int64 `json:"bits,omitempty"`
}

type ChatMessageReply struct {
	ParentMessageId   string `json:"parent_message_id,omitempty"`
	ParentMessageBody string `json:"parent_message_body,omitempty"`
	ParentUserId      string `json:"parent_user_id,omitempty"`
	ParentUserName    string `json:"parent_user_name,omitempty"`
	ParentUserLogin   string `json:"parent_user_login,omitempty"`
	ThreadMessageId   string `json:"thread_message_id,omitempty"`
	ThreadUserId      string `json:"thread_user_id,omitempty"`
	ThreadUserName    string `json:"thread_user_name,omitempty"`
	ThreadUserLogin   string `json:"thread_user_login,omitempty"`
}
