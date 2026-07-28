package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/entities/platform"
)

type User struct {
	ID                uuid.UUID         `json:"id,omitempty"`
	Platform          platform.Platform `json:"platform,omitempty"`
	PlatformID        string            `json:"platform_id,omitempty"`
	TokenID           null.String       `json:"token_id"`
	IsBotAdmin        bool              `json:"is_bot_admin,omitempty"`
	ApiKey            string            `json:"api_key,omitempty"`
	IsBanned          bool              `json:"is_banned,omitempty"`
	HideOnLandingPage bool              `json:"hide_on_landing_page,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	Login             string            `json:"login,omitempty"`
	DisplayName       string            `json:"display_name,omitempty"`
	Avatar            string            `json:"avatar,omitempty"`

	isNil bool
}

func (c User) IsNil() bool {
	return c.isNil
}

var Nil = User{
	isNil: true,
}

type OnlineUser struct {
	ID        uuid.UUID
	ChannelID uuid.UUID
	UserID    uuid.UUID
	UserName  string

	isNil bool
}

func (c OnlineUser) IsNil() bool {
	return c.isNil
}

var NilOnlineUser = OnlineUser{
	isNil: true,
}
