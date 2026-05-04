package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type Channels struct {
	ID               string    `gorm:"primaryKey;column:id;type:UUID;"                json:"id"`
	TwitchUserID     *string   `gorm:"column:twitch_user_id;type:UUID;"               json:"twitchUserId"`
	TwitchBotEnabled bool      `gorm:"column:twitch_bot_enabled;type:BOOL;default:false;" json:"twitchBotEnabled"`
	KickUserID       *string   `gorm:"column:kick_user_id;type:UUID;"                 json:"kickUserId"`
	KickBotEnabled   bool      `gorm:"column:kick_bot_enabled;type:BOOL;default:false;" json:"kickBotEnabled"`
	IsEnabled        bool      `gorm:"column:isEnabled;type:BOOL;"       json:"isEnabled"`
	IsTwitchBanned   bool      `gorm:"column:isTwitchBanned;type:BOOL;" json:"isTwitchBanned"`
	IsBotMod         bool      `gorm:"column:isBotMod;type:BOOL;" json:"isBotMod"`
	BotID            string    `gorm:"column:botId;type:TEXT;"                        json:"botId"`
	PlanID           *string   `gorm:"column:plan_id;type:UUID;"                      json:"planId"`
	CreatedAt        time.Time `gorm:"column:created_at;type:TIMESTAMPTZ;default:now()" json:"createdAt"`

	Commands []ChannelsCommands `gorm:"foreignKey:ChannelID" json:"commands"`
	Roles    []*ChannelRole     `gorm:"foreignKey:ChannelID" json:"roles"`
	User     *Users             `gorm:"foreignKey:TwitchUserID;references:ID" json:"user"`
}

func (c *Channels) TableName() string {
	return "channels"
}

func (c *Channels) Platform() string {
	if c.KickUserID != nil {
		return "kick"
	}
	if c.TwitchUserID != nil {
		return "twitch"
	}
	return ""
}

func (c *Channels) IsOwner(userID string) bool {
	if c.TwitchUserID != nil && *c.TwitchUserID == userID {
		return true
	}
	if c.KickUserID != nil && *c.KickUserID == userID {
		return true
	}
	return false
}
