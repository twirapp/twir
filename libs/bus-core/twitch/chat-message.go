package twitch

import (
	"time"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

type TwitchChatMessage struct {
	Message                     *ChatMessageMessage `json:"message,omitempty"`
	Cheer                       *ChatMessageCheer   `json:"cheer,omitempty"`
	Reply                       *ChatMessageReply   `json:"reply,omitempty"`
	ID                          string              `json:"id"`
	BroadcasterUserId           string              `json:"broadcaster_user_id"`
	BroadcasterUserName         string              `json:"broadcaster_user_name"`
	BroadcasterUserLogin        string              `json:"broadcaster_user_login"`
	ChatterUserId               string              `json:"chatter_user_id"`
	ChatterUserName             string              `json:"chatter_user_name"`
	ChatterUserLogin            string              `json:"chatter_user_login"`
	MessageId                   string              `json:"message_id"`
	Color                       string              `json:"color"`
	MessageType                 string              `json:"message_type"`
	ChannelPointsCustomRewardId string              `json:"channel_points_custom_reward_id"`
	Badges                      []ChatMessageBadge  `json:"badges,omitempty"`

	EnrichedData ChatMessageEnrichedData `json:"enriched_data,omitempty"`
}

type ChatMessageEnrichedData struct {
	UsedEmotesWithThirdParty map[string]int        `json:"used_emotes_with_third_party"`
	ChannelCommandPrefix     string                `json:"channel_command_prefix"`
	DbChannel                channelsmodel.Channel `json:"db_channel"`
	ChannelStream            *streamsmodel.Stream  `json:"channel_stream"`
	DbUser                   *DbUser               `json:"db_user,omitempty"`
	DbUserChannelStat        *DbUserChannelStat    `json:"db_user_channel_stat,omitempty"`
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
	Id         string   `json:"id,omitempty"`
	EmoteSetId string   `json:"emote_set_id,omitempty"`
	OwnerId    string   `json:"owner_id,omitempty"`
	Format     []string `json:"format,omitempty"`
}

type ChatMessageMessageFragmentMention struct {
	UserId    string `json:"user_id,omitempty"`
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

type ChatMessageBadge struct {
	Id    string `json:"id,omitempty"`
	SetId string `json:"set_id,omitempty"`
	Info  string `json:"info,omitempty"`
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
